#! /bin/bash

set -e

watch linkerd --context=east -n test stat \
  --from deploy/linkerd-gateway \
  --from-namespace linkerd-multicluster \
  deploy/podinfo
