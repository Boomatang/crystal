"""
Kubernetes resource templates for generating mock admission review events.

This package provides Python functions to create realistic Kubernetes resource
structures (Deployments, ConfigMaps, Secrets) wrapped in AdmissionReview objects.
"""

from .admission_review import create_admission_review
from .deployment import create_deployment

__all__ = [
    'create_admission_review',
    'create_deployment',
]
