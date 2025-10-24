"""
Secret template for Kubernetes Secret resources.

This is a placeholder template for future Secret event generation.
"""

import base64
import uuid
from datetime import datetime, timezone
from typing import Optional


def create_secret(
    name: str,
    namespace: str = "default",
    data: Optional[dict] = None,
    secret_type: str = "Opaque",
    uid: Optional[str] = None,
    resource_version: str = "12345",
    generation: int = 1,
    creation_timestamp: Optional[str] = None,
    labels: Optional[dict] = None,
) -> dict:
    """
    Create a Kubernetes Secret object structure.

    Args:
        name: The Secret name
        namespace: The Secret namespace (default: "default")
        data: Secret data as key-value pairs (will be base64 encoded)
        secret_type: Secret type (default: "Opaque")
        uid: The Secret UID (auto-generated if None)
        resource_version: The resource version (default: "12345")
        generation: The generation number (default: 1)
        creation_timestamp: ISO8601 timestamp (auto-generated if None)
        labels: Custom labels dict (default: {"app": name})

    Returns:
        dict: Complete Secret object structure

    Note:
        Data values will be automatically base64 encoded if they are strings.
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

    # Base64 encode data values if they're strings
    encoded_data = {}
    for key, value in data.items():
        if isinstance(value, str):
            encoded_data[key] = base64.b64encode(value.encode()).decode()
        else:
            encoded_data[key] = value

    secret = {
        "apiVersion": "v1",
        "kind": "Secret",
        "metadata": {
            "name": name,
            "namespace": namespace,
            "uid": uid,
            "resourceVersion": resource_version,
            "generation": generation,
            "creationTimestamp": creation_timestamp,
            "labels": labels.copy(),
        },
        "type": secret_type,
        "data": encoded_data,
    }

    return secret
