# Crystal Demo Tools

Interactive web interface and mock event generators for the Crystal policy machinery system.

## Overview

This directory contains two main tools:

1. **Streamlit Frontend** (`app.py`): Interactive web interface for visualizing policy workflows and resource topology
2. **Mock Event Generator** (`mock_deployment_updates.py`): Python script to generate realistic Kubernetes admission webhook events

### Frontend Capabilities

The Crystal frontend is a Streamlit-based web application that provides real-time visualization of:

- **Workflow Graphs**: Visual representation of the policy workflow execution flow
- **Node Topology**: Directed acyclic graph (DAG) of Kubernetes resources and their relationships
- **Subgraph Exploration**: Filter and view specific resource relationships by kind, name, and direction

### Mock Event Generator

Generate realistic Kubernetes UPDATE events for testing and development:

- Simulates admission webhook events with incrementing replica counts
- Maintains realistic metadata (resourceVersion, generation, timestamps)
- Extensible template system for adding ConfigMaps, Secrets, etc.

## Prerequisites

- Python 3.13+
- Poetry (Python dependency management)
- Graphviz (graph rendering)
- Crystal backend server running on http://localhost:8000

## Installation

Install dependencies using Poetry:

```sh
poetry install
```

This will install:
- **streamlit**: Web application framework
- **graphviz**: Graph visualization library
- **requests**: HTTP client for backend API calls
- **mkdocs**: Documentation generation (optional)

## Running the Frontend

Start the Streamlit application:

```sh
poetry run streamlit run app.py
```

The frontend will be available at:
- **Frontend UI**: http://localhost:8501

**Note**: Ensure the Crystal backend server is running before starting the frontend.

## Features

### Workflow Graph Visualization

Displays the complete workflow execution flow including:
- Pre-conditions
- Actions (sequential execution)
- Post-conditions
- Error handlers
- Nested workflows (sub-worlds)

Click **"Fetch and Render Workflow Graph"** to retrieve and display the current workflow configuration.

### Node Topology Graph

Shows the complete Kubernetes resource topology:
- All tracked resources (Deployments, ConfigMaps, Secrets, etc.)
- Parent-child relationships between resources
- Resource dependencies established through custom linkers

Click **"Fetch and Render Node Graph"** to view the full resource topology.

### Subgraph Filtering

Explore specific portions of the resource topology:

1. **Kind**: Resource type (e.g., Deployment, ConfigMap, Secret)
2. **Name**: Resource name (e.g., my-app, app-config)
3. **Direction**:
   - `up` - Show only parent resources
   - `down` - Show only child resources
   - `both` - Show full relationship tree (default)

This allows focused analysis of specific resource dependencies and impact analysis.

## API Integration

The frontend connects to the Crystal backend API endpoints:

- `GET /graph` - Workflow graph in DOT format
- `GET /nodelist` - Node topology graph in DOT format
- `GET /nodelist?kind=X&name=Y&direction=Z` - Filtered subgraph

**Backend URL**: Configured in `app.py` (default: `http://localhost:8000`)

## Configuration

### Backend Endpoint

To connect to a different backend server, modify the `endpoint` variable in `app.py`:

```python
endpoint = "http://localhost:8000"  # Change to your backend URL
```

### Streamlit Configuration

Create `.streamlit/config.toml` for custom Streamlit settings:

```toml
[server]
port = 8501
address = "localhost"

[theme]
primaryColor = "#F63366"
backgroundColor = "#FFFFFF"
secondaryBackgroundColor = "#F0F2F6"
textColor = "#262730"
font = "sans serif"
```

## Troubleshooting

### Backend Connection Errors

If you see "Failed to fetch graph" errors:

1. Verify the backend is running:
   ```sh
   curl http://localhost:8000/graph
   ```

2. Check the backend logs for errors:
   ```sh
   LOG_LEVEL=DEBUG go run cmd/example_one/main.go
   ```

3. Ensure no firewall blocking port 8000

### Graphviz Rendering Issues

If graphs don't render:

1. Install system Graphviz package:
   ```sh
   # Ubuntu/Debian
   sudo apt-get install graphviz

   # macOS
   brew install graphviz

   # Fedora/RHEL
   sudo dnf install graphviz
   ```

2. Reinstall Python graphviz package:
   ```sh
   poetry install --sync
   ```

## Development

### Adding New Visualizations

1. Add API client function to `src/crystal/__init__.py`:
   ```python
   def fetch_new_endpoint(url):
       resp = requests.get(f"{url}/new-endpoint", timeout=10)
       if resp.status_code == 200:
           return resp.text
       st.error(f"Failed to fetch data: {resp.status_code}")
   ```

2. Add UI component to `app.py`:
   ```python
   if st.button("Fetch New Data"):
       data = fetch_new_endpoint(endpoint)
       if data:
           st.graphviz_chart(data)
   ```

### Running in Production

For production deployment, use a production-grade server:

```sh
poetry run streamlit run app.py --server.port=8501 --server.address=0.0.0.0
```

## Mock Event Generator Usage

### Quick Start

Generate Deployment UPDATE events with incrementing replica counts:

```sh
cd demo
poetry run example_two
```

Or using the direct Python module invocation:

```sh
cd demo
poetry run python -m crystal.scripts.mock_deployment_updates
```

This will:
1. Send UPDATE events incrementing replicas from 1 to 5
2. Wait 2 seconds between each event
3. POST to `http://localhost:8000/event`

### Command-Line Options

```sh
poetry run example_two --help
```

**Available options:**

- `--backend-url URL` - Backend endpoint (default: `http://localhost:8000/event`)
- `--replicas MIN-MAX` - Replica range to cycle through (default: `1-5`)
- `--sleep SECONDS` - Delay between events (default: `2.0`)
- `--deployment-name NAME` - Deployment name (default: `demo-deployment`)
- `--namespace NS` - Kubernetes namespace (default: `default`)
- `--loop` - Loop continuously, resetting to MIN replicas after reaching MAX
- `--verbose` - Show full request/response JSON bodies

### Example Usage

**Send 10 events with 3-second delays:**
```sh
poetry run example_two --replicas 1-10 --sleep 3
```

**Continuous loop for stress testing:**
```sh
poetry run example_two --replicas 1-5 --sleep 1 --loop
```

**Custom deployment with verbose output:**
```sh
poetry run example_two \
  --deployment-name my-app \
  --namespace production \
  --replicas 3-10 \
  --verbose
```

### Output Example

```
============================================================
Mock Deployment UPDATE Event Generator
============================================================
Backend URL:      http://localhost:8000/event
Deployment:       demo-deployment
Namespace:        default
Replica range:    1 -> 5
Sleep duration:   2.0s
Loop mode:        disabled
============================================================

[10:30:45.123] Sending UPDATE event #1
  Deployment: demo-deployment (replicas: 2)
  Generation: 2, ResourceVersion: 12346
✓ Event sent successfully (200 OK)

[10:30:47.456] Sending UPDATE event #2
  Deployment: demo-deployment (replicas: 3)
  Generation: 3, ResourceVersion: 12347
✓ Event sent successfully (200 OK)
```

## Extending the Mock Generator

### Adding New Resource Types

The template system makes it easy to add support for ConfigMaps, Secrets, and other resources.

**Template structure:**

```
demo/src/crystal/
├── templates/
│   ├── __init__.py              # Export templates
│   ├── admission_review.py      # AdmissionReview wrapper
│   ├── deployment.py            # Deployment template
│   ├── configmap.py             # Future: ConfigMap template
│   └── secret.py                # Future: Secret template
└── scripts/
    ├── __init__.py
    └── mock_deployment_updates.py  # Main event generator
```

### Creating a ConfigMap Template

The ConfigMap template already exists at `src/crystal/templates/configmap.py`!

**Create a new script `src/crystal/scripts/mock_configmap_updates.py`:**

```python
from crystal.templates.admission_review import create_admission_review
from crystal.templates.configmap import create_configmap

# Generate ConfigMap events
config = create_configmap(
    name="app-config",
    namespace="default",
    data={"key": "value"}
)

event = create_admission_review(
    operation="UPDATE",
    resource_object=config,
    kind="ConfigMap",
    group="",  # Core API group
    version="v1",
    name="app-config",
    namespace="default"
)

# POST event...
```

### Template Function Reference

**`create_admission_review()`** - Wrap resources in AdmissionReview

Parameters:
- `operation` - CREATE, UPDATE, or DELETE
- `resource_object` - The Kubernetes resource dict
- `kind` - Resource kind (Deployment, ConfigMap, etc.)
- `group` - API group ("apps", "" for core)
- `version` - API version ("v1")
- `name` - Resource name
- `namespace` - Resource namespace
- `request_uid` - Unique request ID (auto-generated)
- `request_timestamp` - Request timestamp (auto-generated)
- `old_object` - Previous state for UPDATE (optional)

**`create_deployment()`** - Generate Deployment objects

Parameters:
- `name` - Deployment name
- `namespace` - Namespace (default: "default")
- `replicas` - Replica count (default: 1)
- `uid` - Resource UID (auto-generated)
- `resource_version` - ResourceVersion (default: "12345")
- `generation` - Generation number (default: 1)
- `creation_timestamp` - Creation time (auto-generated)
- `labels` - Custom labels dict
- `image` - Container image (default: "nginx:latest")
- `container_name` - Container name (default: same as deployment)
- `env_from_configmap` - ConfigMap name for env vars

## License

This project is a proof-of-concept and is not intended for production use.
