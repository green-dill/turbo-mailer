import { exec } from 'child_process';
import { promisify } from 'util';
import { writeFileSync } from 'fs';

const execAsync = promisify(exec);

function isInSubnet(ip: string, cidr: string): boolean {
  const [subnet, bits] = cidr.split('/');
  const mask = ~(2 ** (32 - parseInt(bits)) - 1);
  const ipNum = ip.split('.').reduce((acc, octet) => (acc << 8) + parseInt(octet), 0) >>> 0;
  const subnetNum = subnet.split('.').reduce((acc, octet) => (acc << 8) + parseInt(octet), 0) >>> 0;
  return (ipNum & mask) === (subnetNum & mask);
}

async function runCommand(command: string): Promise<string> {
  try {
    const { stdout, stderr } = await execAsync(command);
    if (stderr) console.error(stderr);
    return stdout.trim();
  } catch (error) {
    console.error(`Error executing command: ${command}`);
    console.error(error);
    throw error;
  }
}

async function installPackages(): Promise<void> {
  const aptPackages = [
    'curl',
    'wget',
    'vim',
    'make',
    'git',
    'haproxy',
  ];

  console.log('Updating package lists...');
  await runCommand('apt-get update');

  console.log('Installing packages...');
  await runCommand(`apt-get install -yq ${aptPackages.join(' ')}`);

  console.log('Installing kubectl...');
  await runCommand(`
    curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/
  `);

  console.log('Installing kustomize...');
  await runCommand(`
    curl -s "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh" | bash && \
    mv kustomize /usr/local/bin/
  `);

  console.log('Installing Helm...');
  await runCommand('curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash');

  console.log('Installing Helmfile...');
  const helmfileVersion = '0.169.0';
  await runCommand(`
    curl -Lo helmfile.tar.gz https://github.com/helmfile/helmfile/releases/download/v${helmfileVersion}/helmfile_${helmfileVersion}_linux_amd64.tar.gz && \
    tar -xzf helmfile.tar.gz && \
    chmod +x helmfile && \
    mv helmfile /usr/local/bin/ && \
    rm helmfile.tar.gz
  `);
}

async function getInternalIP(): Promise<string> {
  const ipOutput = await runCommand('ip -4 addr show');
  const lines = ipOutput.split('\n');
  const internalInterfaces = ['eth', 'ens', 'enp', 'eno'];
  const internalCIDRs = [
    '10.0.0.0/8',
    '172.16.0.0/12',
    '192.168.0.0/16'
  ];

  let internalIP: string | null = null;
  let fallbackIP: string | null = null;

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (internalInterfaces.some(iface => line.includes(iface))) {
      const ipLine = lines[i + 1];
      if (ipLine) {
        const match = ipLine.match(/inet (\d+\.\d+\.\d+\.\d+)/);
        if (match) {
          const ip = match[1];
          if (internalCIDRs.some(cidr => isInSubnet(ip, cidr))) {
            return ip;
          } else if (!fallbackIP) {
            fallbackIP = ip;
          }
        }
      }
    }
  }

  if (fallbackIP) {
    console.warn('Warning: No internal IP found. Using fallback IP:', fallbackIP);
    return fallbackIP;
  }

  throw new Error('Unable to determine IP address');
}

async function installK3s(): Promise<void> {
  console.log('Installing K3s...');

  try {
    const k3sStatus = await runCommand('systemctl status k3s');
    if (k3sStatus.includes('Active: active')) {
      console.log('K3s is already installed and running. Skipping installation...');
      return;
    }
  } catch (error) {
    console.log('K3s is not installed. Proceeding with installation...');
  }

  const k3sVersion = "v1.30.1+k3s1";
  const k3sToken = "09ffc3c3f9badaa9297";
  const internalIP = await getInternalIP();

  const k3sInstallScript = `
    curl --retry 10 --retry-delay 0 -sfL https://get.k3s.io | \
    sed 's/curl -o/curl -C - --retry 10 --retry-delay 0 --connect-timeout 5 -o/g' | \
    INSTALL_K3S_VERSION="${k3sVersion}" \
    INSTALL_K3S_EXEC="--disable traefik,servicelb" \
    sh -s - \
     --token ${k3sToken} \
     --cluster-init \
     --cluster-cidr=10.62.8.0/22 \
     --service-cidr=10.62.12.0/22 \
     --kube-controller-manager-arg=node-cidr-mask-size=25 \
     --embedded-registry
  `;
  await runCommand(k3sInstallScript);

  console.log('Enabling K3s service for automatic startup...');
  await runCommand('systemctl enable k3s --now');

  console.log('Configuring K3s registries...');
  const registriesConfig = `
mirrors:
  "*":
`;
  writeFileSync('/etc/rancher/k3s/registries.yaml', registriesConfig);

  // Generate agent installation script
  const agentInstallScript = `#!/bin/bash
curl --retry 10 --retry-delay 0 --retry-all-errors -sfL https://get.k3s.io | \
  sed 's/curl -o/curl -C - --retry 10 --retry-delay 0 --connect-timeout 5 --retry-all-errors -o/g' | \
  INSTALL_K3S_VERSION="${k3sVersion}" \
  sh -s - \
   agent --server https://${internalIP}:6443 \
   --token ${k3sToken}
`;

  writeFileSync('install-agent.sh', agentInstallScript, { mode: 0o755 });

  console.log('\nK3s installation completed.');
  console.log('To add an agent to this cluster, run the following command on the agent node:');
  console.log(agentInstallScript);
  console.log('\nAn install-agent.sh script has been created in the current directory.');
}

async function checkK3sInstallation(): Promise<void> {
  console.log('Checking K3s installation...');
  let retries = 100;
  const retryInterval = 10000; // 10 seconds

  while (retries > 0) {
    try {
      const nodes = await runCommand('k3s kubectl get node');
      if (nodes.includes('Ready')) {
        console.log('K3s is installed and running correctly.');
        return;
      }
    } catch (error) {
      console.log(`K3s is not ready yet. Retrying in ${retryInterval / 1000} seconds...`);
    }

    await new Promise(resolve => setTimeout(resolve, retryInterval));
    retries--;
  }

  throw new Error('K3s installation check failed after multiple retries.');
}

async function installBaseServices(): Promise<void> {
  console.log('Installing base services...');
  await runCommand('make sync cert-manager');
  await runCommand('make sync ingress-nginx');

  console.log('Applying kustomizations...');
  process.chdir('kustomize');
  await runCommand('kubectl apply -k .');
}

async function installApplications(): Promise<void> {
  console.log('Installing applications...');
  await runCommand('helmfile -e prod sync');
}

async function main(): Promise<void> {
  console.log('Starting installation process...');
  const workingDir = process.cwd();

  await installPackages();
  await installK3s();

  await checkK3sInstallation();

  process.chdir(workingDir + '/../infra/cluster');
  await installBaseServices();
  process.chdir(workingDir + '/../');
  await installApplications();

  console.log('Installation process completed successfully!');
}

main().catch(error => {
  console.error('An error occurred during the installation process:');
  console.error(error);
  process.exit(1);
});
