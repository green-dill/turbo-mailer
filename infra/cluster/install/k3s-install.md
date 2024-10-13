# k3s install

## Update and install basic tools

```bash
apt update
apt install -y curl wget vim
```

## Install k3s master

```bash
curl --retry 10 --retry-delay 0 -sfL https://get.k3s.io | \
  sed 's/curl -o/curl -C - --retry 10 --retry-delay 0 --connect-timeout 5 -o/g' | \
  INSTALL_K3S_VERSION="v1.30.1+k3s1" \
  INSTALL_K3S_EXEC="--disable traefik,servicelb" \
  sh -s - \
   --token 09ffc3c3f9badaa9297 \
   --cluster-init \
   --cluster-cidr=10.62.8.0/22 \
   --service-cidr=10.62.12.0/22 \
   --kube-controller-manager-arg=node-cidr-mask-size=25 \
   --embedded-registry

cat <<'EOF' > /etc/rancher/k3s/registries.yaml
mirrors:
  "*":
EOF
```

## Join k3s node

> replace xxx with your master node ip

```bash
curl --retry 10 --retry-delay 0 --retry-all-errors -sfL https://get.k3s.io | \
  sed 's/curl -o/curl -C - --retry 10 --retry-delay 0 --connect-timeout 5 --retry-all-errors -o/g' | \
  INSTALL_K3S_VERSION="v1.30.1+k3s1" \
  sh -s - \
   agent --server https://xxx:6443 \
   --token 09ffc3c3f9badaa9297
```
