#! /bin/bash

set -e

for ctx in west east; do
  echo "Adding test services on cluster: ${ctx} ........."
  kubectl --context=${ctx} apply \
    -k "github.com/linkerd/website/multicluster/${ctx}/"
  kubectl --context=${ctx} -n test \
    rollout status deploy/podinfo || break
  echo "-------------\n"
done

kubectl --context=west -n test port-forward svc/frontend 8080 &

sleep 3

curl http://localhost:8080

open http://localhost:8080