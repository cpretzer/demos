#!/bin/bash

# check for linkerd and kubectl
if ! type kubectl > /dev/null; then
    echo "\n\nPlease install kubectl:"
    echo "\n\nhttps://kubernetes.io/docs/tasks/tools/#kubectl"
    exit 1
fi

if ! type linkerd > /dev/null; then
    echo "\n\nInstall the linkerd cli with the command:"
    echo "\n\n curl -sL https://run.linked.io/install | sh"
    exit 1
fi


# install linkerd

# install linkerd-viz

# deploy emojivoto

# install linkerd-failover

# deploy failover YAML
