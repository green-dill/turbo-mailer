# Turbo Mailer

![Kubernetes](https://img.shields.io/badge/kubernetes-%23326ce5.svg?style=flat&logo=kubernetes&logoColor=white)
![Helm](https://img.shields.io/badge/helm-%230F1689.svg?style=flat&logo=helm&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)
![TypeScript](https://img.shields.io/badge/typescript-%23007ACC.svg?style=flat&logo=typescript&logoColor=white)
![React](https://img.shields.io/badge/react-%2320232a.svg?style=flat&logo=react&logoColor=%2361DAFB)
![CI](https://github.com/green-dill/turbo-mailer/actions/workflows/deploy.yaml/badge.svg?branch=develop)

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


