#! /bin/bash

linkerd install | kubectl apply -f -

linkerd check

# Install Flagger
kubectl apply -k github.com/weaveworks/flagger/kustomize/linkerd

kubectl -n linkerd rollout status deploy/flagger --watch

kubectl create ns test 
kubectl annotate ns test linkerd.io/inject=enabled
kubectl apply -f slow-cooker.yml
kubectl apply -k github.com/weaveworks/flagger/kustomize/podinfo

kubectl rollout -n test status deploy podinfo --watch

kubectl apply -n test -f frontend.yml

kubectl port-forward  -n test svc/frontend 8080

kubectl apply -f canary.yml

kubectl -n canary-deployment get ev --watch

kubectl -n test set image deployment/podinfo \
  podinfod=stefanprodan/podinfo:3.1.1

kubectl -n test get ev --watch

kubectl -n test get canary --watch

kubectl -n test get trafficsplit podinfo -o yaml

watch linkerd stat ts -n test

watch linkerd -n test stat deploy --from deploy/load

kubectl -n test port-forward svc/frontend 8080