#! /bin/bash

set -e

linkerd --context=west -n test stat --from deploy/frontend svc

linkerd dashboard &