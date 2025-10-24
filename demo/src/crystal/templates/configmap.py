"""
ConfigMap template for Kubernetes ConfigMap resources.

This is a placeholder template for future ConfigMap event generation.
"""

import uuid
from datetime import datetime, timezone
from typing import Optional


def create_configmap(
    name: str,
    namespace: str = "default",
    data: Optional[dict] = None,
    uid: Optional[str] = None,
    resource_version: str = "12345",
    generation: int = 1,
    creation_timestamp: Optional[str] = None,
    labels: Optional[dict] = None,
) -> dict:
    """
    Create a Kubernetes ConfigMap object structure.

    Args:
        name: The ConfigMap name
        namespace: The ConfigMap namespace (default: "default")
        data: ConfigMap data as key-value pairs (default: {})
        uid: The ConfigMap UID (auto-generated if None)
        resource_version: The resource version (default: "12345")
        generation: The generation number (default: 1)
        creation_timestamp: ISO8601 timestamp (auto-generated if None)
        labels: Custom labels dict (default: {"app": name})

    Returns:
        dict: Complete ConfigMap object structure
    """
    if uid is None:
        uid = str(uuid.uuid4())

    if creation_timestamp is None:
        now = datetime.now(timezone.utc)
        creation_timestamp = now.strftime("%Y-%m-%dT%H:%M:%SZ")

    if labels is None:
        labels = {"app": name}

    if data is None:
        data = {}

    configmap = {
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "name": name,
            "namespace": namespace,
            "uid": uid,
            "resourceVersion": resource_version,
            "generation": generation,
            "creationTimestamp": creation_timestamp,
            "labels": labels.copy(),
        },
        "data": data.copy(),
    }

    return configmap
