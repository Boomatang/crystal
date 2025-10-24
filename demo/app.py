import streamlit as st

from crystal import fetch_dot_graph, fetch_node_graph


def main():
    endpoint = "http://localhost:8000"

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


if __name__ == "__main__":
    main()
