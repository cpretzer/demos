#! /bin/bash 

set -e

for ctx in west east; do
  echo "Checking cluster: ${ctx} .........\n"
  linkerd --context=${ctx} check  --pre
  echo "-------------\n"
done
