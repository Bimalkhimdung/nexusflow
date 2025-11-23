# NexusFlow

**The best self-hosted, open-source, modern alternative to Jira**

[![License: Apache 2.0](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://go.dev/)
[![Kubernetes](https://img.shields.io/badge/Kubernetes-Ready-326CE5?logo=kubernetes)](https://kubernetes.io/)

NexusFlow is a 100% headless, API-first project management platform built for startups, scale-ups, and DevOps teams (10-500 users). Designed to be cheaper, faster, and more extensible than Atlassian products.

## ✨ Features

### MVP (v1.0)

- 🔐 **OAuth2/OIDC Authentication** via Ory Hydra + Kratos
- 👥 **Organization & Team Management** with invite links
- 📊 **Unlimited Projects** with Kanban, Scrum, and Bug-tracking templates
- 🎯 **Issue Hierarchy** - Epics → Stories → Sub-tasks
- 🔧 **Custom Fields** - 10+ field types with full flexibility
- 🔄 **Visual Workflow Designer** - Drag-drop statuses, transitions, and rules
- 📋 **Kanban Boards** - WIP limits, swimlanes, filters, real-time drag-drop
- 🏃 **Scrum Support** - Backlog, sprints, goals, burndown charts
- 📈 **Dashboards** - Personal & project dashboards with 10+ gadgets
- 💬 **Rich Collaboration** - Comments, @mentions, reactions, real-time updates
- 📎 **Attachments** - File uploads with thumbnails and previews
- 🔍 **Full-text Search** - Elasticsearch-powered with JQL-like syntax
- 🔔 **Notifications** - In-app, email, Slack, webhooks, WebSocket
- 🔗 **Git Integrations** - GitHub, GitLab, Bitbucket commit & PR linking
- 🚀 **Complete APIs** - REST + GraphQL + WebSocket + 50+ webhook events
- 🔒 **RBAC** - Role-based permissions with project-level roles

### Post-MVP (v1.1+)

- SAML & social login
- Advanced reports & velocity charts
- Automation rules engine
- Cross-project roadmaps
- Time tracking & worklogs
- Mobile apps (React Native)

## 🏗️ Architecture

NexusFlow is built as a modern microservices architecture:

```
┌─────────────────────────────────────────────────────────────┐
│                      Traefik Gateway                        │
│              (API Gateway + Load Balancer)                  │
└─────────────────────────────────────────────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        │                     │                     │
        ▼                     ▼                     ▼
┌──────────────┐      ┌──────────────┐     ┌──────────────┐
│ Auth Service │      │ User Service │     │  Org Service │
│ (Ory Hydra)  │      │              │     │              │
└──────────────┘      └──────────────┘     └──────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
                    ┌─────────┴─────────┐
                    │   Apache Kafka    │
                    │  (Event Streaming) │
                    └───────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐      ┌──────────────┐     ┌──────────────┐
│Project Svc   │      │ Issue Service│     │Workflow Svc  │
└──────────────┘      └──────────────┘     └──────────────┘
        │                     │                     │
        └─────────────────────┼─────────────────────┘
                              │
        ┌─────────────────────┼─────────────────────┐
        ▼                     ▼                     ▼
┌──────────────┐      ┌──────────────┐     ┌──────────────┐
│ PostgreSQL   │      │Elasticsearch │     │    MinIO     │
│  (Primary)   │      │   (Search)   │     │  (Storage)   │
└──────────────┘      └──────────────┘     └──────────────┘
```

### Tech Stack

- **Backend**: Go 1.24+ for all microservices
- **Message Broker**: Apache Kafka (Redpanda compatible)
- **API Gateway**: Traefik v3+
- **Authentication**: Ory Hydra (OAuth2/OIDC) + Ory Kratos
- **Database**: PostgreSQL 16+ (multi-tenant)
- **Search**: Elasticsearch 8.x / OpenSearch
- **Storage**: MinIO (S3-compatible)
- **Frontend**: TypeScript + React 19 + Vite + TanStack
- **Real-time**: WebSocket + Server-Sent Events
- **Deployment**: Docker + Kubernetes + Helm + ArgoCD
- **Observability**: OpenTelemetry, Prometheus, Grafana, Loki, Jaeger

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- kubectl (for Kubernetes deployment)
- Helm 3+ (for Kubernetes deployment)

### Local Development

1. **Clone the repository**

   ```bash
   git clone https://github.com/yourusername/nexusflow.git
   cd nexusflow
   ```

2. **Start infrastructure services**

   ```bash
   docker-compose up -d
   ```

3. **Generate protobuf code**

   ```bash
   make generate-proto
   ```

4. **Run all services**

   ```bash
   make run-all
   ```

5. **Access the application**
   - API Gateway: <http://localhost:8080>
   - Traefik Dashboard: <http://localhost:8081>
   - MinIO Console: <http://localhost:9001>

### Kubernetes Deployment

1. **Install with Helm**

   ```bash
   helm install nexusflow ./deployments/helm/nexusflow
   ```

2. **Access the application**

   ```bash
   kubectl port-forward svc/traefik 8080:80
   ```

For detailed setup instructions, see [Development Setup Guide](docs/development/setup.md).

## 📚 Documentation

- [Architecture Overview](docs/architecture/overview.md)
- [Microservices Documentation](docs/architecture/microservices.md)
- [Development Setup](docs/development/setup.md)
- [API Documentation](docs/api/README.md)
- [Coding Conventions](docs/development/conventions.md)

## 🛠️ Development

### Project Structure

```
nexusflow/
├── services/           # All microservices
│   ├── user-service/
│   ├── org-service/
│   ├── project-service/
│   └── ...
├── pkg/               # Shared libraries
│   ├── logger/
│   ├── config/
│   ├── database/
│   └── kafka/
├── proto/             # Protobuf definitions
│   ├── common/
│   ├── user/
│   ├── org/
│   └── ...
├── deployments/       # Infrastructure as Code
│   ├── helm/
│   └── docker-compose.yml
├── docs/              # Documentation
└── scripts/           # Build and utility scripts
```

### Available Make Commands

```bash
make generate-proto    # Generate Go code from protobuf
make build-all        # Build all services
make test-all         # Run all tests
make lint             # Run linters
make run-all          # Run all services locally
make docker-build     # Build Docker images
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run tests (`make test-all`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## 📊 Project Goals

### 12 Months After Launch

- 15,000+ GitHub stars
- 1,000+ production clusters
- Average 50-user cluster cost < $120/month
- Community contributes >25% of commits

## 📄 License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## 🙏 Acknowledgments

Built with:

- [Ory](https://www.ory.sh/) - Authentication & authorization
- [Traefik](https://traefik.io/) - API Gateway
- [Apache Kafka](https://kafka.apache.org/) - Event streaming
- [PostgreSQL](https://www.postgresql.org/) - Primary database
- [Elasticsearch](https://www.elastic.co/) - Search engine
- [MinIO](https://min.io/) - Object storage

## 📞 Support

- 📖 [Documentation](docs/)
- 💬 [Discussions](https://github.com/yourusername/nexusflow/discussions)
- 🐛 [Issue Tracker](https://github.com/yourusername/nexusflow/issues)
- 💼 [Commercial Support](https://nexusflow.io/support)

---

**Made with ❤️ by the NexusFlow team**
