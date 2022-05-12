#!/bin/bash
ORG_DOMAIN="${ORG_DOMAIN:-k3d.example.com}"
LINKERD="${LINKERD:-linkerd}"

# check for linkerd, k3d, step, helm, and kubectl
if ! type kubectl > /dev/null; then
    echo "\n\nPlease install kubectl:"
    echo "\n\nhttps://kubernetes.io/docs/tasks/tools/#kubectl"
    exit 1
fi

if ! type linkerd > /dev/null; then
    echo "\n\nInstall the linkerd cli with the command:"
    echo "\n\n curl -sL https://run.linked.io/install | sh"
    exit 1
fi

if ! type k3d > /dev/null; then
    echo "\n\nInstall k3d:"
    echo "\n\n https://k3d.io/v5.4.1/#installation"
    exit 1
fi

if ! type step > /dev/null; then
    echo "\n\nInstall step:"
    echo "\n\n https://smallstep.com/docs/step-cli/installation"
    exit 1
fi

if ! type helm > /dev/null; then
    echo "\n\nInstall helm:"
    echo "\n\n https://helm.sh/docs/intro/install/"
    exit 1
fi

# create the clusters
port=6440
for cluster in dev east west ; do
    if k3d cluster get "$cluster" >/dev/null 2>&1 ; then
        echo "Already exists: $cluster" >&2
    else
        k3d cluster create "$cluster" \
            --api-port="$((port++))" \
            --network=multicluster-example \
            --k3s-server-arg="--cluster-domain=$cluster.${ORG_DOMAIN}" \
            --wait
    fi
done

# install linkerd


# Generate the trust roots. These never touch the cluster. In the real world
# we'd squirrel these away in a vault.
step certificate create \
    "identity.linkerd.${ORG_DOMAIN}" \
    ca.crt ca.key \
    --profile root-ca \
    --no-password  --insecure --force

for cluster in dev east west ; do
    # Check that the cluster is up and running.
    while ! $LINKERD --context="k3d-$cluster" check --pre ; do :; done

    # Create issuing credentials. These end up on the cluster (and can be
    # rotated from the root).
    crt="${cluster}-issuer.crt"
    key="${cluster}-issuer.key"
    domain="${cluster}.${ORG_DOMAIN}"
    step certificate create "identity.linkerd.${domain}" \
        "$crt" "$key" \
        --ca=ca.crt \
        --ca-key=ca.key \
        --profile=intermediate-ca \
        --not-after 8760h --no-password --insecure

    # Install Linkerd into the cluster.
    $LINKERD --context="k3d-$cluster" install \
            --proxy-log-level="linkerd=debug,trust_dns=debug,info" \
            --cluster-domain="$domain" \
            --identity-trust-domain="$domain" \
            --identity-trust-anchors-file=ca.crt \
            --identity-issuer-certificate-file="${crt}" \
            --identity-issuer-key-file="${key}" |
        kubectl --context="k3d-$cluster" apply -f -

    # Wait some time and check that the cluster has started properly.
    sleep 30
    while ! $LINKERD --context="k3d-$cluster" check ; do :; done

    $LINKERD --context="k3d-$cluster" viz install |
        kubectl --context="k3d-$cluster" apply -f -

    # Setup the multicluster components on the server
    $LINKERD --context="k3d-$cluster" multicluster install |
        kubectl --context="k3d-$cluster" apply -f -
done

# install emojivoto
for cluster in dev east west ; do
    kubectl --context="k3d-$cluster" apply -f https://run.linkerd.io/emojivoto.yml
done

# install linkerd-failover
# In case you haven't added the linkerd-edge repo already
helm repo add linkerd-edge https://helm.linkerd.io/edge
helm repo up

for cluster in dev east west ; do
    helm install linkerd-failover -n linkerd-failover --create-namespace \
    --kube-context $context --devel linkerd-edge/linkerd-failover
done

# deploy failover YAML
kubectl apply -f failover-config
