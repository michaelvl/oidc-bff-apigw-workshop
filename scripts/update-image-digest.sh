#! /bin/bash

set -e

TAG=$1

SCRIPT=$(readlink -f $0)
SCRIPTPATH=`dirname $SCRIPT`

if [ -z "$TAG" ]; then
    SHA=`git rev-parse --short HEAD`
    TAG="sha-$SHA"
    echo "No image tag specified, using HEAD: $TAG"
fi

IMAGE=ghcr.io/michaelvl/oidc-bff-apigw-workshop
DIGEST=$($SCRIPTPATH/../scripts/skopeo.sh inspect docker://$IMAGE:$TAG | jq -r .Digest)
echo "Using digest: $DIGEST"
sed -i -E "s#(.*?ghcr.io/michaelvl/oidc-bff-apigw-workshop@).*#\1$DIGEST#" kubernetes/spa-cdn.yaml
sed -i -E "s#(.*?ghcr.io/michaelvl/oidc-bff-apigw-workshop@).*#\1$DIGEST#" kubernetes/spa-api-gw.yaml
sed -i -E "s#(.*?ghcr.io/michaelvl/oidc-bff-apigw-workshop@).*#\1$DIGEST#" kubernetes/protected-api.yaml
