#! /bin/bash

set -e

k3d cluster create l5d-policy --wait

ledge install --registry ghcr.io/linkerd \
  -f config.yml |
  kubectl apply -f -

# HACK: sleep before calling wait
sleep 20

# Wait for linkerd control plane
kubectl wait -n linkerd \
  --for=condition=ready pod \
  --selector=linkerd.io/control-plane-ns=linkerd \
  --timeout=240s

# Install linkerd viz
ledge viz install | kubectl apply -f -

# HACK: sleep before calling wait
sleep 5

ledge inject https://run.linkerd.io/emojivoto.yml | kubectl apply -f -

kubectl wait -n linkerd-viz \
  --for=condition=ready pod \
  --selector=linkerd.io/extension=viz \
  --timeout=120s

ledge viz stat po -n emojivoto

# kubectl apply -f https://run.linkerd.io/emojivoto-policy.yml

