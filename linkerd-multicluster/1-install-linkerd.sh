#! /bin/bash 

set -e

linkerd install \
  --identity-trust-anchors-file new-certs/root.crt \
  --identity-issuer-certificate-file new-certs/issuer.crt \
  --identity-issuer-key-file new-certs/issuer.key \
  | tee \
    >(kubectl --context=west apply -f -) \
    >(kubectl --context=east apply -f -)

for ctx in west east; do
  echo "Checking cluster: ${ctx} .........\n"
  linkerd --context=${ctx} check 
  echo "-------------\n"
done