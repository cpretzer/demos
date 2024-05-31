#! /bin/bash

set -e

linkerd --context=west -n test tap deploy/frontend | \
  grep "$(kubectl --context=east -n linkerd-multicluster get svc linkerd-gateway \
    -o "custom-columns=GATEWAY_IP:.status.loadBalancer.ingress[*].ip")"

kubectl --context=west -n test run -it --rm --image=alpine:3 test -- \
  /bin/sh -c "apk add curl && curl -vv http://podinfo-east:9898"