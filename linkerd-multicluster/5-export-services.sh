#! /bin/bash

set -e

kubectl --context=east get svc -n test podinfo -o yaml | \
  linkerd multicluster export-service - | \
  kubectl --context=east apply -f -

kubectl --context=west -n test get svc podinfo-east