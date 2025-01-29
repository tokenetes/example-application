#!/bin/bash

info() {
    printf "\e[34m%s\n\e[0m" "$1"
}

error() {
    printf "\e[31m%s\n\e[0m" "$1"
}

success() {
    printf "\e[32m%s\n\e[0m" "$1"
}

resource_exists() {
    resource_type=$1
    resource_name=$2
    kubectl get ${resource_type} ${resource_name} > /dev/null 2>&1
}

delete_k8s_resource() {
    resource_type=$1
    resource_name=$2
    if resource_exists ${resource_type} ${resource_name}; then
        kubectl delete ${resource_type} ${resource_name} || { error "Failed to delete ${resource_type} ${resource_name}"; exit 1; }
    else
        info "${resource_type} ${resource_name} does not exist, skipping deletion."
    fi
}

info "Uninstalling tconfigd..."

delete_k8s_resource namespace tokenetes-system
delete_k8s_resource clusterrole tconfigd-service-account-role
delete_k8s_resource clusterrolebinding tconfigd-service-account-binding
delete_k8s_resource mutatingwebhookconfiguration tokenetes-agent-injector
delete_k8s_resource crd trats.tokenetes.io
delete_k8s_resource crd tokenetesconfigs.tokenetes.io

success "tconfigd uninstalled successfully."
