import requests
import streamlit as st


def send_event(url, deployment_name, namespace, resource_version):
    """
    Sends an AdmissionReview event to the server.

    Parameters:
        url: Base server URL (e.g., "http://localhost:8000")
        deployment_name: Name of the deployment resource
        namespace: Kubernetes namespace
        resource_version: Integer resourceVersion for the event

    Returns:
        Response dict from server
    """
    event = {
        "kind": "AdmissionReview",
        "apiVersion": "admission.k8s.io/v1",
        "request": {
            "uid": f"{deployment_name}-{resource_version}",
            "kind": {"group": "apps", "version": "v1", "kind": "Deployment"},
            "resource": {"group": "apps", "version": "v1", "resource": "deployments"},
            "namespace": namespace,
            "name": deployment_name,
            "operation": "UPDATE",
            "object": {
                "apiVersion": "apps/v1",
                "kind": "Deployment",
                "metadata": {
                    "name": deployment_name,
                    "namespace": namespace,
                    "resourceVersion": str(resource_version),
                },
            },
        },
    }

    resp = requests.post(f"{url}/event", json=event, timeout=10)
    if resp.status_code in (200, 201):
        return resp.json()

    st.error(f"Failed to send event: {resp.status_code}")
    return None


def fetch_event_list(url):
    """
    Fetches current queue state from /list endpoint.

    Parameters:
        url: Base server URL

    Returns:
        List of AdmissionReview events currently in queue
    """
    resp = requests.get(f"{url}/list", timeout=10)
    if resp.status_code == 200:
        return resp.json()

    st.error(f"Failed to fetch event list: {resp.status_code}")
    return []


def fetch_dot_graph(url):
    resp = requests.get(f"{url}/graph", timeout=10)
    if resp.status_code == 200:
        return resp.text

    st.error(f"Failed to fetch graph: {resp.status_code}")


def fetch_node_graph(url, kind=None, name=None, direction=None):
    params = {}
    if kind is not None and name is not None:
        params = {"kind": kind, "name": name, "direction": direction}

    resp = requests.get(f"{url}/nodelist", params=params, timeout=10)
    if resp.status_code == 200:
        return resp.text

    st.error(f"Failed to fetch graph: {resp.status_code}")
