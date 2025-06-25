import streamlit as st

from crystal import fetch_dot_graph


def main():
    endpoint = "http://localhost:8000"

    st.title("Render workflow")
    if st.button("Fetch and Render Workflow Graph"):
        data = fetch_dot_graph(endpoint)
        if data:
            st.graphviz_chart(data)


if __name__ == "__main__":
    main()
