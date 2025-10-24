import requests
import streamlit as st


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
