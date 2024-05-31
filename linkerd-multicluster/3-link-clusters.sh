#! /bin/bash

set -e 

linkerd --context=east multicluster link --cluster-name east |
  kubectl --context=west apply -f -

linkerd --context=west check --multicluster

linkerd --context=west multicluster gateways
