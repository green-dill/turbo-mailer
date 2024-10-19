import { spawn, ChildProcess } from 'child_process';
import inquirer from 'inquirer';

process.env.KUBECONFIG = `${process.cwd()}/kubeconfig-readonly-dev.yaml`;

const namespace = 'turbo-mailer-dev';

const portForwards = [
  { name: 'PostgreSQL', command: `kubectl port-forward -n ${namespace} service/postgresql --address 0.0.0.0 15432:5432` },
  { name: 'Redis', command: `kubectl port-forward -n ${namespace} service/redis-master --address 0.0.0.0 16379:6379` },
  { name: 'Turbo Mailer API', command: `kubectl port-forward -n ${namespace} service/turbo-mailer-api --address 0.0.0.0 5000:5000` },
  { name: 'Postfix SMTP', command: `kubectl port-forward -n ${namespace} service/postfix-smtp --address 0.0.0.0 2525:25` },
];

const childProcesses: ChildProcess[] = [];

function startPortForward(service: { name: string; command: string }) {
  const [cmd, ...args] = service.command.split(' ');
  const process = spawn(cmd, args);
  childProcesses.push(process);

  process.stdout.on('data', (data) => {
    console.log(`[${service.name}] stdout: ${data.toString().trimEnd()}`);
  });

  process.stderr.on('data', (data) => {
    console.error(`[${service.name}] stderr: ${data.toString().trimEnd()}`);
  });

  process.on('close', (code) => {
    console.log(`[${service.name}] Child process exited with code ${code}`);
  });
}

function cleanupProcesses() {
  console.log('Cleaning up child processes...');
  childProcesses.forEach((process) => {
    if (!process.killed) {
      process.kill();
    }
  });
}

function main() {
  inquirer
    .prompt([
      {
        type: 'checkbox',
        name: 'selectedServices',
        message: 'Select services to port-forward:',
        choices: portForwards.map(service => ({
          name: service.name,
          value: service,
          // default selected services
          checked: ['PostgreSQL', 'Redis'].includes(service.name)
        })),
      }
    ])
    .then((answers: { selectedServices: typeof portForwards }) => {
      const selectedServices = answers.selectedServices;
      if (selectedServices.length === 0) {
        console.log('No services selected. Exiting...');
        process.exit(0);
      }

      console.log('Starting port forwarding for selected services...');
      selectedServices.forEach(startPortForward);
    });

  // Handle SIGINT (Ctrl+C) and SIGTERM
  process.on('SIGINT', () => {
    console.log('Received SIGINT. Cleaning up...');
    cleanupProcesses();
    process.exit(0);
  });

  process.on('SIGTERM', () => {
    console.log('Received SIGTERM. Cleaning up...');
    cleanupProcesses();
    process.exit(0);
  });
}

main();
