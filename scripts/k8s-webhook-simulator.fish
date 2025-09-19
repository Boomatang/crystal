#!/usr/bin/env fish

# Kubernetes Webhook Simulator for Controller-Runtime
# Simulates the HTTP data sent by Kubernetes API server to operator controllers

set -l LOCALHOST_URL "http://localhost:8080"
set -l WEBHOOK_PATH "/mutate"
set -l VALIDATE_PATH "/validate"

function send_configmap_create_webhook
    echo "Sending ConfigMap CREATE webhook..."
    # Simulate AdmissionReview for ConfigMap creation
     echo -n '{
            "apiVersion": "admission.k8s.io/v1",
            "kind": "AdmissionReview",
            "request": {
                "uid": "705ab4f5-6393-11e8-b7cc-42010a800002",
                "kind": {
                    "group": "",
                    "version": "v1",
                    "kind": "ConfigMap"
                },
                "resource": {
                    "group": "",
                    "version": "v1",
                    "resource": "configmaps"
                },
                "requestKind": {
                    "group": "",
                    "version": "v1",
                    "kind": "ConfigMap"
                },
                "requestResource": {
                    "group": "",
                    "version": "v1",
                    "resource": "configmaps"
                },
                "name": "exampleconfig",
                "namespace": "default",
                "operation": "CREATE",
                "userInfo": {
                    "username": "admin",
                    "uid": "admin",
                    "groups": ["system:masters"]
                },
                "object": {
                    "apiVersion": "v1",
                    "kind": "ConfigMap",
                    "metadata": {
                        "name": "exampleconfig",
                        "namespace": "default"
                    },
                    "data": {
                        "app.properties": "debug=true\nlog.level=info",
                        "database.url": "postgresql://localhost:5432/mydb"
                    }
                },
                "oldObject": null,
                "dryRun": false,
                "options": {
                    "apiVersion": "meta.k8s.io/v1",
                    "kind": "CreateOptions"
                }
            }
    }' | http -v POST :8000/event \
    "Content-Type: application/json" \
        "Accept: application/json" 
end

function send_configmap_update_webhook
    echo "Sending ConfigMap UPDATE webhook..."
    
    http POST ":8080" \
        "Content-Type: application/json" \
        "Accept: application/json" \
        body:='{
            "apiVersion": "admission.k8s.io/v1",
            "kind": "AdmissionReview",
            "request": {
                "uid": "705ab4f5-6393-11e8-b7cc-42010a800003",
                "kind": {
                    "group": "",
                    "version": "v1",
                    "kind": "ConfigMap"
                },
                "resource": {
                    "group": "",
                    "version": "v1",
                    "resource": "configmaps"
                },
                "requestKind": {
                    "group": "",
                    "version": "v1",
                    "kind": "ConfigMap"
                },
                "requestResource": {
                    "group": "",
                    "version": "v1",
                    "resource": "configmaps"
                },
                "name": "example-config",
                "namespace": "default",
                "operation": "UPDATE",
                "userInfo": {
                    "username": "admin",
                    "uid": "admin",
                    "groups": ["system:masters"]
                },
                "object": {
                    "apiVersion": "v1",
                    "kind": "ConfigMap",
                    "metadata": {
                        "name": "example-config",
                        "namespace": "default",
                        "resourceVersion": "12345"
                    },
                    "data": {
                        "app.properties": "debug=false\nlog.level=warn",
                        "database.url": "postgresql://localhost:5432/mydb",
                        "cache.enabled": "true"
                    }
                },
                "oldObject": {
                    "apiVersion": "v1",
                    "kind": "ConfigMap",
                    "metadata": {
                        "name": "example-config",
                        "namespace": "default",
                        "resourceVersion": "12344"
                    },
                    "data": {
                        "app.properties": "debug=true\nlog.level=info",
                        "database.url": "postgresql://localhost:5432/mydb"
                    }
                },
                "dryRun": false,
                "options": {
                    "apiVersion": "meta.k8s.io/v1",
                    "kind": "UpdateOptions"
                }
            }
        }'
end

function send_configmap_delete_webhook
    echo "Sending ConfigMap DELETE webhook..."
    
    http POST "$LOCALHOST_URL$WEBHOOK_PATH" \
        "Content-Type: application/json" \
        "Accept: application/json" \
        body:='{
            "apiVersion": "admission.k8s.io/v1",
            "kind": "AdmissionReview",
            "request": {
                "uid": "705ab4f5-6393-11e8-b7cc-42010a800004",
                "kind": {
                    "group": "",
                    "version": "v1",
                    "kind": "ConfigMap"
                },
                "resource": {
                    "group": "",
                    "version": "v1",
                    "resource": "configmaps"
                },
                "requestKind": {
                    "group": "",
                    "version": "v1",
                    "kind": "ConfigMap"
                },
                "requestResource": {
                    "group": "",
                    "version": "v1",
                    "resource": "configmaps"
                },
                "name": "example-config",
                "namespace": "default",
                "operation": "DELETE",
                "userInfo": {
                    "username": "admin",
                    "uid": "admin",
                    "groups": ["system:masters"]
                },
                "object": null,
                "oldObject": {
                    "apiVersion": "v1",
                    "kind": "ConfigMap",
                    "metadata": {
                        "name": "example-config",
                        "namespace": "default",
                        "resourceVersion": "12345"
                    },
                    "data": {
                        "app.properties": "debug=false\nlog.level=warn",
                        "database.url": "postgresql://localhost:5432/mydb",
                        "cache.enabled": "true"
                    }
                },
                "dryRun": false,
                "options": {
                    "apiVersion": "meta.k8s.io/v1",
                    "kind": "DeleteOptions"
                }
            }
        }'
end

function send_validation_webhook
    echo "Sending validation webhook..."
    
    http POST "$LOCALHOST_URL$VALIDATE_PATH" \
        "Content-Type: application/json" \
        "Accept: application/json" \
        body:='{
            "apiVersion": "admission.k8s.io/v1",
            "kind": "AdmissionReview",
            "request": {
                "uid": "705ab4f5-6393-11e8-b7cc-42010a800005",
                "kind": {
                    "group": "",
                    "version": "v1",
                    "kind": "ConfigMap"
                },
                "resource": {
                    "group": "",
                    "version": "v1",
                    "resource": "configmaps"
                },
                "requestKind": {
                    "group": "",
                    "version": "v1",
                    "kind": "ConfigMap"
                },
                "requestResource": {
                    "group": "",
                    "version": "v1",
                    "resource": "configmaps"
                },
                "name": "example-config",
                "namespace": "default",
                "operation": "CREATE",
                "userInfo": {
                    "username": "admin",
                    "uid": "admin",
                    "groups": ["system:masters"]
                },
                "object": {
                    "apiVersion": "v1",
                    "kind": "ConfigMap",
                    "metadata": {
                        "name": "example-config",
                        "namespace": "default"
                    },
                    "data": {
                        "app.properties": "debug=true\nlog.level=info",
                        "database.url": "postgresql://localhost:5432/mydb"
                    }
                },
                "oldObject": null,
                "dryRun": false,
                "options": {
                    "apiVersion": "meta.k8s.io/v1",
                    "kind": "CreateOptions"
                }
            }
        }'
end

function run_all_webhooks
    echo "Running all Kubernetes webhook simulations..."
    echo "=========================================="
    
    send_configmap_create_webhook
    echo ""
    sleep 1
    
    send_configmap_update_webhook
    echo ""
    sleep 1
    
    send_configmap_delete_webhook
    echo ""
    sleep 1
    
    send_validation_webhook
    echo ""
    
    echo "All webhooks sent successfully!"
end

# Main execution
if test (count $argv) -eq 0
    run_all_webhooks
else
    switch $argv[1]
        case "create"
            send_configmap_create_webhook
        case "update"
            send_configmap_update_webhook
        case "delete"
            send_configmap_delete_webhook
        case "validate"
            send_validation_webhook
        case "all"
            run_all_webhooks
        case "*"
            echo "Usage: $argv[0] [create|update|delete|validate|all]"
            echo "  create   - Send ConfigMap CREATE webhook"
            echo "  update   - Send ConfigMap UPDATE webhook"
            echo "  delete   - Send ConfigMap DELETE webhook"
            echo "  validate - Send validation webhook"
            echo "  all      - Send all webhooks (default)"
    end
end 
