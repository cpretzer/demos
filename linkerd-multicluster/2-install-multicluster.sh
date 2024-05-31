#! /bin/bash

set -e

for ctx in west east; do
  echo "Installing on cluster: ${ctx} ........."
  linkerd --context=${ctx} multicluster install | \
    kubectl --context=${ctx} apply -f - || break
  echo "-------------\n"
done

for ctx in west east; do
  echo "Checking gateway on cluster: ${ctx} ........."
  kubectl --context=${ctx} -n linkerd-multicluster \
    rollout status deploy/linkerd-gateway || break
  echo "-------------\n"
done

for ctx in west east; do
  echo "Checking cluster: ${ctx} ........."
  linkerd --context=${ctx} check --multicluster || break
  echo "-------------\n"
done

kubectl --context=west -n linkerd-multicluster get all