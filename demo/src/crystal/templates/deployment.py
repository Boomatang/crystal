"""
Deployment template for Kubernetes Deployment resources.
"""

import uuid
from datetime import datetime, timezone
from typing import Optional


def create_deployment(
    name: str,
    namespace: str = "default",
    replicas: int = 1,
    uid: Optional[str] = None,
    resource_version: str = "12345",
    generation: int = 1,
    creation_timestamp: Optional[str] = None,
    labels: Optional[dict] = None,
    image: str = "nginx:latest",
    container_name: Optional[str] = None,
    env_from_configmap: Optional[str] = None,
) -> dict:
    """
    Create a Kubernetes Deployment object structure.

    Args:
        name: The deployment name
        namespace: The deployment namespace (default: "default")
        replicas: Number of replicas (default: 1)
        uid: The deployment UID (auto-generated if None)
        resource_version: The resource version (default: "12345")
        generation: The generation number (default: 1)
        creation_timestamp: ISO8601 timestamp (auto-generated if None)
        labels: Custom labels dict (default: {"app": name})
        image: Container image (default: "nginx:latest")
        container_name: Container name (default: same as deployment name)
        env_from_configmap: ConfigMap name to load env vars from (optional)

    Returns:
        dict: Complete Deployment object structure
    """
    if uid is None:
        # Generate k8s-style UUID
        uid = str(uuid.uuid4())

    if creation_timestamp is None:
        # Generate timestamp without microseconds: "2025-10-24T10:30:45Z"
        now = datetime.now(timezone.utc)
        creation_timestamp = now.strftime("%Y-%m-%dT%H:%M:%SZ")

    if labels is None:
        labels = {"app": name}

    if container_name is None:
        container_name = name

    # Build container spec
    container_spec = {
        "name": container_name,
        "image": image,
    }

    # Add envFrom if configmap is specified
    if env_from_configmap:
        container_spec["envFrom"] = [
            {
                "configMapRef": {
                    "name": env_from_configmap
                }
            }
        ]

    deployment = {
        "apiVersion": "apps/v1",
        "kind": "Deployment",
        "metadata": {
            "name": name,
            "namespace": namespace,
            "uid": uid,
            "resourceVersion": resource_version,
            "generation": generation,
            "creationTimestamp": creation_timestamp,
            "labels": labels.copy(),  # Copy to avoid mutation
        },
        "spec": {
            "replicas": replicas,
            "selector": {
                "matchLabels": {
                    "app": name
                }
            },
            "template": {
                "metadata": {
                    "labels": {
                        "app": name
                    }
                },
                "spec": {
                    "containers": [container_spec]
                }
            }
        }
    }

    return deployment
