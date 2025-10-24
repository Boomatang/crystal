#!/usr/bin/env python3
"""
Mock Kubernetes Deployment UPDATE event generator.

This script generates realistic Kubernetes admission webhook UPDATE events
for a Deployment resource, incrementing replica counts and metadata fields
with each iteration.

Usage:
    poetry run python mock_deployment_updates.py --replicas 1-5 --sleep 3
"""

import argparse
import sys
import time
import uuid
from datetime import datetime, timezone
from typing import Tuple

import requests

from crystal.templates.admission_review import create_admission_review
from crystal.templates.deployment import create_deployment


class DeploymentState:
    """Maintains state for a deployment across multiple events."""

    def __init__(
        self,
        name: str,
        namespace: str,
        initial_replicas: int,
        initial_resource_version: int = 12345,
    ):
        self.name = name
        self.namespace = namespace
        self.uid = str(uuid.uuid4())
        self.creation_timestamp = datetime.now(timezone.utc).strftime(
            "%Y-%m-%dT%H:%M:%SZ"
        )
        self.resource_version = initial_resource_version
        self.generation = 1
        self.replicas = initial_replicas

    def increment_for_update(self, new_replicas: int) -> None:
        """Update state for a new UPDATE event."""
        self.replicas = new_replicas
        self.resource_version += 1
        self.generation += 1


def parse_replica_range(range_str: str) -> Tuple[int, int]:
    """
    Parse replica range string like '1-5' into (min, max).

    Args:
        range_str: Range string in format 'MIN-MAX'

    Returns:
        Tuple of (min_replicas, max_replicas)

    Raises:
        ValueError: If range format is invalid
    """
    try:
        parts = range_str.split("-")
        if len(parts) != 2:
            raise ValueError("Range must be in format 'MIN-MAX'")
        min_val, max_val = int(parts[0]), int(parts[1])
        if min_val >= max_val:
            raise ValueError("MIN must be less than MAX")
        if min_val < 0:
            raise ValueError("MIN must be non-negative")
        return min_val, max_val
    except (ValueError, IndexError) as e:
        raise ValueError(f"Invalid replica range '{range_str}': {e}")


def send_update_event(
    backend_url: str,
    state: DeploymentState,
    event_num: int,
    verbose: bool = False,
) -> bool:
    """
    Send an UPDATE event to the backend.

    Args:
        backend_url: Backend URL to POST to
        state: Current deployment state
        event_num: Event number for logging
        verbose: Show full request/response

    Returns:
        True if successful, False otherwise
    """
    # Create deployment object
    deployment = create_deployment(
        name=state.name,
        namespace=state.namespace,
        replicas=state.replicas,
        uid=state.uid,
        resource_version=str(state.resource_version),
        generation=state.generation,
        creation_timestamp=state.creation_timestamp,
        labels={"app": state.name},
        image="nginx:latest",
    )

    # Wrap in AdmissionReview
    admission_review = create_admission_review(
        operation="UPDATE",
        resource_object=deployment,
        kind="Deployment",
        group="apps",
        version="v1",
        name=state.name,
        namespace=state.namespace,
    )

    # Print event info
    timestamp = datetime.now().strftime("%H:%M:%S.%f")[:-3]
    print(f"\n[{timestamp}] Sending UPDATE event #{event_num}")
    print(f"  Deployment: {state.name} (replicas: {state.replicas})")
    print(
        f"  Generation: {state.generation}, ResourceVersion: {state.resource_version}"
    )

    if verbose:
        import json

        print("\n--- Request Body ---")
        print(json.dumps(admission_review, indent=2))

    # Send request
    try:
        response = requests.post(
            backend_url,
            json=admission_review,
            headers={"Content-Type": "application/json"},
            timeout=5,
        )
        response.raise_for_status()

        print(
            f"\u2713 Event sent successfully ({response.status_code} {response.reason})"
        )

        if verbose:
            print("\n--- Response ---")
            print(f"Status: {response.status_code}")
            print(f"Body: {response.text}")

        return True

    except requests.exceptions.RequestException as e:
        print(f"\u2717 Failed to send event: {e}")
        return False


def main():
    """Main entry point."""
    parser = argparse.ArgumentParser(
        description="Generate mock Kubernetes Deployment UPDATE events",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument(
        "--backend-url",
        default="http://localhost:8000/event",
        help="Backend URL to POST events to",
    )
    parser.add_argument(
        "--replicas",
        default="1-5",
        help="Replica range in format MIN-MAX (e.g., 1-5, 3-10)",
    )
    parser.add_argument(
        "--sleep",
        type=float,
        default=2.0,
        help="Seconds to sleep between events",
    )
    parser.add_argument(
        "--deployment-name",
        default="demoDeployment",
        help="Deployment name",
    )
    parser.add_argument(
        "--namespace",
        default="default",
        help="Deployment namespace",
    )
    parser.add_argument(
        "--loop",
        action="store_true",
        help="Loop back to start after reaching max replicas",
    )
    parser.add_argument(
        "--verbose",
        action="store_true",
        help="Show full request/response bodies",
    )

    args = parser.parse_args()

    # Parse replica range
    try:
        min_replicas, max_replicas = parse_replica_range(args.replicas)
    except ValueError as e:
        print(f"Error: {e}", file=sys.stderr)
        sys.exit(1)

    print("=" * 60)
    print("Mock Deployment UPDATE Event Generator")
    print("=" * 60)
    print(f"Backend URL:      {args.backend_url}")
    print(f"Deployment:       {args.deployment_name}")
    print(f"Namespace:        {args.namespace}")
    print(f"Replica range:    {min_replicas} -> {max_replicas}")
    print(f"Sleep duration:   {args.sleep}s")
    print(f"Loop mode:        {'enabled' if args.loop else 'disabled'}")
    print("=" * 60)

    # Initialize state
    state = DeploymentState(
        name=args.deployment_name,
        namespace=args.namespace,
        initial_replicas=min_replicas,
    )

    event_num = 1

    try:
        while True:
            # Generate events for each replica count
            for replicas in range(min_replicas + 1, max_replicas + 1):
                state.increment_for_update(replicas)

                success = send_update_event(
                    backend_url=args.backend_url,
                    state=state,
                    event_num=event_num,
                    verbose=args.verbose,
                )

                if not success:
                    print("\nWarning: Failed to send event, continuing...")

                event_num += 1

                # Sleep before next event (except after last event if not looping)
                if args.loop or replicas < max_replicas:
                    time.sleep(args.sleep)

            # Exit if not looping
            if not args.loop:
                print(f"\n\u2713 Completed {event_num - 1} events")
                break

            # Reset replica count for next loop
            print("\n--- Looping back to start ---")
            state.replicas = min_replicas

    except KeyboardInterrupt:
        print(f"\n\nInterrupted. Sent {event_num - 1} events.")
        sys.exit(0)


if __name__ == "__main__":
    main()
