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
go run cmd/example_one/main.go
```

**With custom logging configuration:**
```sh
# Set log level (DEBUG, INFO, WARN, ERROR)
LOG_LEVEL=DEBUG go run cmd/example_one/main.go

# Use JSON format for structured logging
LOG_FORMAT=json go run cmd/example_one/main.go

# Combine both
LOG_LEVEL=DEBUG LOG_FORMAT=json go run cmd/example_one/main.go
```

**Frontend Interface:**
```sh
cd frontend
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
- **Structured Logging**: Configurable log levels and formats (text/JSON) for better debugging and production use

## Configuration

### Logging

Crystal uses structured logging with configurable levels and formats:

**Environment Variables:**
- `LOG_LEVEL`: Set logging verbosity (DEBUG, INFO, WARN, ERROR). Default: INFO
- `LOG_FORMAT`: Output format (text, json). Default: text

**Log Levels:**
- **DEBUG**: Detailed execution flow, useful for development and troubleshooting
- **INFO**: Normal operational messages (workflow started, nodes added, etc.)
- **WARN**: Unexpected but recoverable conditions
- **ERROR**: Failures that need attention

**Examples:**
```sh
# Development with detailed logs
LOG_LEVEL=DEBUG go run cmd/example_one/main.go

# Production with JSON logs for log aggregation
LOG_LEVEL=INFO LOG_FORMAT=json go run cmd/example_one/main.go

# Minimal logging for performance
LOG_LEVEL=ERROR go run cmd/example_one/main.go
```

## API Endpoints

- `GET /metrics` - Prometheus metrics
- `GET /graph` - Workflow graph (DOT format)
- `GET /nodelist` - Node relationship graph (DOT format)
- `POST /event` - Webhook event handler
- `GET /list` - List all events


## License

This project is a proof-of-concept and is not intended for production use.
