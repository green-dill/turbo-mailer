# Turbo Mailer

![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=flat&logo=kubernetes&logoColor=white)
![Helm](https://img.shields.io/badge/helm-%230F1689.svg?style=flat&logo=helm&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)
![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=flat&logo=typescript&logoColor=white)
![React](https://img.shields.io/badge/react-%2320232a.svg?style=flat&logo=react&logoColor=%2361DAFB)
![CI](https://github.com/green-dill/turbo-mailer/actions/workflows/deploy.yaml/badge.svg?branch=develop)
![CI](https://github.com/green-dill/turbo-mailer/actions/workflows/postfix.yaml/badge.svg?branch=develop)

<img src="docs/assets/banner.png" alt="Turbo Mailer" height="20%">

Turbo Mailer is a high-performance marketing email delivery system, built with Golang for the backend and TypeScript/React for the frontend. It is deployed using Kubernetes (k8s) and Helm, offering a powerful, cost-effective, and easy-to-maintain solution for businesses of all sizes.

## Key Features

- **High Performance**: Leveraging Golang's concurrency model, Turbo Mailer can handle millions of emails per hour with minimal resource usage.
- **Cost-Effective**: Designed to run efficiently on Kubernetes, reducing infrastructure costs while maintaining high availability.
- **Easy Maintenance**: Built with DevOps best practices, making it simple to deploy, update, and monitor.
- **Scalable**: Automatically scales based on demand, ensuring optimal performance during peak times.
- **Robust Monitoring**: Integrated with popular monitoring tools for real-time insights into system health and performance.
- **Customizable**: Flexible architecture allows for easy integration with existing systems and customization of email templates.

## Deployment

### Kubernetes-based Deployment

This assumes you already have a Kubernetes cluster (cloud-provider managed or self-hosted). For self-hosted clusters, we recommend using k3s. You can refer to the [k3s-install.md](infra/cluster/install/k3s-install.md) document for guidance.

### Application Deployment

By default, deployment is automated using GitHub Actions. Manual deployment is also an option.

To deploy to the production environment, run the following command:

```bash
helmfile -e prod sync
```

> For specific parameters, refer to `deploy/values-prod.yaml`

### Core/Infrastructure Services

Core/infrastructure services primarily include cluster monitoring, log collection, certificate management, ingress, etc. Install these as needed.

For integrated services, check the [infra](infra) directory.

Here's an example of installing cert-manager:

```bash
cd infra/cluster

make sync cert-manager
```

> You may need to adjust the parameters in `values.yaml` which includes ingress domain, etc.

## One-click Deployment

We provide a one-click deployment script for easy setup. This script will install the necessary dependencies, clone the repository, and run the installation script.

requirements:

- A server running Ubuntu or Debian
- Root access to the server
- `curl` command installed (usually pre-installed on most systems)

if not, install it first:
```bash
apt update && apt install -y curl
```

```bash
curl -sSL https://yinheli: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@raw.githubusercontent.com/green-dill/turbo-mailer/master/install/install.sh | bash
```

> [!WARNING]
> The PAT will expire on 2025-10-01. If expired, please contact maintainer to update the PAT.

> Note: While the one-click installation script sets up the basic infrastructure, additional configuration is required before Turbo Mailer is fully operational. This includes setting up domain names, configuring email settings, and other environment-specific parameters. Please refer to the configuration guide for detailed instructions on how to complete the setup after installation.

## Contributors

We would like to thank all the contributors who have helped make Turbo Mailer better:

<div style="display: flex; flex-wrap: wrap; gap: 10px;">
  <a href="https://github.com/yinheli" style="text-decoration: none;">
    <img src="https://github.com/yinheli.png" width="40" height="40" alt="Yinheli" style="border-radius: 50%;">
  </a>
  <a href="https://github.com/grubylee" style="text-decoration: none;">
    <img src="https://github.com/grubylee.png" width="40" height="40" alt="Grubylee" style="border-radius: 50%;">
  </a>
</div>

## Contributing

We welcome contributions from the community! If you'd like to contribute to this project, please read our [Contributing Guidelines](CONTRIBUTING.md) for detailed information on how to get started, our code style, commit message conventions, and the pull request process.

Here's a quick overview of how to contribute:

1. Clone the repository
2. Create a new branch for your work (use `bugfix/`, `feature/`, or `hotfix/` prefixes)
3. Make your changes
4. Commit your changes with a descriptive commit message
5. Push your branch and open a Pull Request

For more detailed information, please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) file.

## Important Notice

**Turbo Mailer is not an open-source project.**

Please be aware of the following points:

- **Proprietary Software**: Turbo Mailer is a proprietary software product. All rights are reserved.
- **No Public Distribution**: The source code, binaries, and related materials are not available for public distribution or use.
- **No Unauthorized Modifications**: Users are not permitted to modify, reverse engineer, or create derivative works based on Turbo Mailer without explicit permission.
- **Confidentiality**: All information related to Turbo Mailer's architecture, algorithms, and implementation details should be treated as confidential.
- **Support**: Technical support is provided exclusively to licensed users as per the terms of their agreement.
- **Updates and Maintenance**: Software updates and maintenance are managed solely by our development team.
- **Feedback**: While we appreciate user feedback, any suggestions or ideas submitted become the property of Turbo Mailer.
For inquiries about licensing, support, or partnership opportunities, please contact our technical support team at yinheli@gmail.com.
