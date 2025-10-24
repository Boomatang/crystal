"""
AdmissionReview template for Kubernetes admission webhook events.
"""

import uuid
from datetime import datetime, timezone
from typing import Optional


def create_admission_review(
    operation: str,
    resource_object: dict,
    kind: str,
    group: str,
    version: str,
    name: str,
    namespace: str = "default",
    request_uid: Optional[str] = None,
    request_timestamp: Optional[str] = None,
    old_object: Optional[dict] = None,
    username: str = "admin",
    user_uid: str = "admin",
    user_groups: Optional[list] = None,
) -> dict:
    """
    Create a Kubernetes AdmissionReview structure.

    Args:
        operation: The operation type (CREATE, UPDATE, DELETE)
        resource_object: The Kubernetes resource object (e.g., Deployment, ConfigMap)
        kind: The resource kind (e.g., "Deployment", "ConfigMap")
        group: The API group (e.g., "apps", "" for core resources)
        version: The API version (e.g., "v1")
        name: The resource name
        namespace: The resource namespace (default: "default")
        request_uid: Unique identifier for this admission request (auto-generated if None)
        request_timestamp: ISO8601 timestamp with microseconds (auto-generated if None)
        old_object: The previous state for UPDATE operations (None for CREATE)
        username: The user making the request (default: "admin")
        user_uid: The user's UID (default: "admin")
        user_groups: List of user groups (default: ["system:masters"])

    Returns:
        dict: Complete AdmissionReview structure
    """
    if request_uid is None:
        request_uid = str(uuid.uuid4())

    if request_timestamp is None:
        # Generate timestamp with microseconds: "2025-10-24T10:30:45.123456Z"
        now = datetime.now(timezone.utc)
        request_timestamp = now.strftime("%Y-%m-%dT%H:%M:%S.%fZ")

    if user_groups is None:
        user_groups = ["system:masters"]

    # Determine API version for the resource
    if group:
        api_version = f"{group}/{version}"
        resource_plural = f"{kind.lower()}s"
    else:
        api_version = version
        resource_plural = f"{kind.lower()}s"

    # Determine options kind based on operation
    options_kind = f"{operation.capitalize()}Options"

    return {
        "apiVersion": "admission.k8s.io/v1",
        "kind": "AdmissionReview",
        "request": {
            "uid": request_uid,
            "kind": {
                "group": group,
                "version": version,
                "kind": kind,
            },
            "resource": {
                "group": group,
                "version": version,
                "resource": resource_plural,
            },
            "requestKind": {
                "group": group,
                "version": version,
                "kind": kind,
            },
            "requestResource": {
                "group": group,
                "version": version,
                "resource": resource_plural,
            },
            "name": name,
            "namespace": namespace,
            "operation": operation,
            "requestReceivedTimestamp": request_timestamp,
            "userInfo": {
                "username": username,
                "uid": user_uid,
                "groups": user_groups,
            },
            "object": resource_object,
            "oldObject": old_object,
            "dryRun": False,
            "options": {
                "apiVersion": "meta.k8s.io/v1",
                "kind": options_kind,
            },
        },
    }
