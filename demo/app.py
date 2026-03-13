import random
import time

import streamlit as st

from crystal import fetch_dot_graph, fetch_event_list, fetch_node_graph, send_event


def get_resource_version_from_event(event):
    """Extract resourceVersion from an AdmissionReview event."""
    try:
        return int(event["request"]["object"]["metadata"]["resourceVersion"])
    except (KeyError, TypeError, ValueError):
        return None


def get_event_key(event):
    """Generate the key for an event (kind/namespace/name)."""
    try:
        req = event["request"]
        return f"{req['kind']['kind']}/{req['namespace']}/{req['name']}"
    except (KeyError, TypeError):
        return None


def deduplication_demo_page(endpoint):
    """Deduplication demo page showing EventQueue behavior."""
    st.title("Event Deduplication Demo")
    st.markdown(
        "Demonstrates how the EventQueue deduplicates events based on resourceVersion."
    )

    # Sidebar controls
    st.sidebar.header("Controls")

    mode = st.sidebar.selectbox("Mode", ["Sequential", "Out-of-Order"])
    deployment_name = st.sidebar.text_input("Deployment Name", value="demo-app")
    namespace = st.sidebar.text_input("Namespace", value="default")
    num_events = st.sidebar.slider("Number of Events", min_value=5, max_value=20, value=10)
    delay_ms = st.sidebar.slider("Delay Between Events (ms)", min_value=0, max_value=500, value=50)

    # Initialize session state
    if "event_log" not in st.session_state:
        st.session_state.event_log = []
    if "stats" not in st.session_state:
        st.session_state.stats = {"sent": 0, "added": 0, "replaced": 0, "skipped": 0}

    # Clear button
    if st.sidebar.button("Clear Log"):
        st.session_state.event_log = []
        st.session_state.stats = {"sent": 0, "added": 0, "replaced": 0, "skipped": 0}
        st.rerun()

    # Send Burst button
    if st.sidebar.button("Send Burst"):
        # Generate resourceVersions based on mode
        base_rv = 100
        if mode == "Sequential":
            versions = list(range(base_rv, base_rv + num_events))
        else:  # Out-of-Order
            versions = list(range(base_rv, base_rv + num_events))
            random.shuffle(versions)

        # Send events
        for rv in versions:
            # Get queue state before
            queue_before = fetch_event_list(endpoint)
            queue_before_map = {}
            for evt in queue_before:
                key = get_event_key(evt)
                if key:
                    queue_before_map[key] = get_resource_version_from_event(evt)

            # Send the event
            send_event(endpoint, deployment_name, namespace, rv)
            st.session_state.stats["sent"] += 1

            # Get queue state after
            queue_after = fetch_event_list(endpoint)
            queue_after_map = {}
            for evt in queue_after:
                key = get_event_key(evt)
                if key:
                    queue_after_map[key] = get_resource_version_from_event(evt)

            # Infer status
            event_key = f"Deployment/{namespace}/{deployment_name}"
            if event_key not in queue_before_map and event_key in queue_after_map:
                status = "Added"
                st.session_state.stats["added"] += 1
            elif (
                event_key in queue_before_map
                and event_key in queue_after_map
                and queue_after_map[event_key] == rv
                and queue_before_map[event_key] < rv
            ):
                status = "Replaced"
                st.session_state.stats["replaced"] += 1
            else:
                status = "Skipped"
                st.session_state.stats["skipped"] += 1

            # Log the event
            st.session_state.event_log.append(
                {
                    "timestamp": time.strftime("%H:%M:%S"),
                    "resourceVersion": rv,
                    "status": status,
                }
            )

            # Delay between events
            if delay_ms > 0:
                time.sleep(delay_ms / 1000.0)

        st.rerun()

    # Main content area
    col1, col2 = st.columns(2)

    with col1:
        st.subheader("Event Log")
        if st.session_state.event_log:
            for entry in st.session_state.event_log:
                status_color = {
                    "Added": "🟢",
                    "Replaced": "🟡",
                    "Skipped": "🔴",
                }.get(entry["status"], "⚪")
                st.text(
                    f"{entry['timestamp']} | rv={entry['resourceVersion']:>3} | {status_color} {entry['status']}"
                )
        else:
            st.info("No events sent yet. Click 'Send Burst' to start.")

    with col2:
        st.subheader("Queue State")
        queue = fetch_event_list(endpoint)
        if queue:
            for evt in queue:
                key = get_event_key(evt)
                rv = get_resource_version_from_event(evt)
                st.text(f"{key}: rv={rv}")
        else:
            st.info("Queue is empty.")

    # Summary stats
    st.subheader("Summary Stats")
    stats = st.session_state.stats
    col_a, col_b, col_c, col_d, col_e = st.columns(5)
    col_a.metric("Events Sent", stats["sent"])
    col_b.metric("Added", stats["added"])
    col_c.metric("Replaced", stats["replaced"])
    col_d.metric("Skipped", stats["skipped"])
    col_e.metric("Queue Size", len(fetch_event_list(endpoint)))


def workflow_page(endpoint):
    """Original workflow visualization page."""
    st.title("Render workflow")

    # Workflow Graph Section
    if st.button("Fetch and Render Workflow Graph"):
        data = fetch_dot_graph(endpoint)
        if data:
            st.graphviz_chart(data)

    # Node Graph Section
    if st.button("Fetch and Render Node Graph"):
        data = fetch_node_graph(endpoint)
        if data:
            st.graphviz_chart(data)

    # Subgraph Section
    st.subheader("Fetch Subgraph")
    col1, col2, col3 = st.columns(3)

    with col1:
        kind = st.text_input("Kind", placeholder="e.g., Deployment, ConfigMap")

    with col2:
        name = st.text_input("Name", placeholder="e.g., my_app")

    with col3:
        direction = st.selectbox("Direction", ["up", "down", "both"], index=2)

    if st.button("Fetch and Render Subgraph"):
        if kind and name:
            data = fetch_node_graph(endpoint, kind=kind, name=name, direction=direction)
            if data:
                st.graphviz_chart(data)
        else:
            st.warning("Please provide both Kind and Name to fetch the subgraph.")


def main():
    endpoint = "http://localhost:8000"

    # Page navigation
    page = st.sidebar.radio("Page", ["Workflow", "Deduplication Demo"])

    if page == "Workflow":
        workflow_page(endpoint)
    else:
        deduplication_demo_page(endpoint)


if __name__ == "__main__":
    main()
