# Crystal - Policy Machinery POC

A proof-of-concept implementation of policy machinery that combines a Go backend with Kubernetes webhook integration and a Python frontend for visualization.

## Overview

Crystal is designed as a microservices-based policy machinery system with a clear separation between the backend processing engine and the frontend visualization layer. It provides:

- **Event-driven workflow system** with graph visualization
- **Kubernetes webhook integration** for resource monitoring
- **Prometheus metrics** for observability
- **Interactive web interface** for system management

## Architecture

- **Backend (Go)**: HTTP server with workflow engine, node management, and webhook processing
- **Frontend (Python)**: Streamlit web interface with graph visualization
- **Data Flow**: Kubernetes webhooks → Event processing → Node creation → Workflow execution → Visualization

## Quick Start

### Prerequisites

- Go 1.22+
- Python 3.13+
- Poetry
- Graphviz

### Running the System

**Backend Server:**
```sh
go run main.go
```

**Frontend Interface:**
```sh
poetry run streamlit run app.py
```

**Access the System:**
- Backend API: http://localhost:8000
- Frontend UI: http://localhost:8501
- Metrics: http://localhost:8000/metrics

## Features

- **Workflow Engine**: Define and execute complex workflows with pre/post conditions
- **Node Management**: Track and link Kubernetes resources automatically
- **Graph Visualization**: Interactive visualization of workflows and resource relationships
- **Webhook Integration**: Process Kubernetes admission reviews in real-time
- **Metrics & Monitoring**: Prometheus metrics for system observability

## API Endpoints

- `GET /metrics` - Prometheus metrics
- `GET /graph` - Workflow graph (DOT format)
- `GET /nodelist` - Node relationship graph (DOT format)
- `POST /event` - Webhook event handler
- `GET /list` - List all events

## Documentation

Comprehensive documentation is available in the `/docs/` directory:

- [Architecture Overview](docs/architecture/overview.md)
- [API Documentation](docs/api/endpoints.md)
- [Workflow System](docs/workflow/overview.md)
- [Development Setup](docs/development/setup.md)
- [Deployment Guide](docs/deployment/installation.md)

## Development

See [Development Setup](docs/development/setup.md) for detailed instructions on setting up a development environment.

## Contributing

See [Contributing Guide](docs/development/contributing.md) for information on contributing to the project.

## License

This project is a proof-of-concept and is not intended for production use.
