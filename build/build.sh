#!/bin/bash

# BUILD_MODE="debugBuild"

BUILD_MODE="minimalBuild"
# DOCKERFILE=alpine.Dockerfile
DOCKERFILE=minimal.Dockerfile

NOTARY_BRANCH=master
NOTARYPKG=github.com/theupdateframework/notary

NOTARY_WRAPPER_BRANCH=master
NOTARY_WRAPPER_PKG=github.com/k8s-gadgets/notary-wrapper

REGISTRY=dgeiger
NOTARY_WRAPPER_TAG=0.0.1

# build notary-wrapper binary
if [ $BUILD_MODE == "debugBuild" ]; then

    docker build \
        --build-arg NOTARY_BRANCH=$NOTARY_BRANCH \
        --build-arg NOTARYPKG=$NOTARYPKG \
        -f ./debug.notary.Dockerfile \
        -t notary-binary \
        .

    docker build \
        --build-arg NOTARY_WRAPPER_BRANCH=$NOTARY_WRAPPER_BRANCH \
        --build-arg NOTARY_WRAPPER_PKG=$NOTARY_WRAPPER_PKG \
        -f ./debug.wrapper.Dockerfile \
        -t $REGISTRY/notary-wrapper:$NOTARY_WRAPPER_TAG \
        --no-cache .

elif [ $BUILD_MODE == "minimalBuild" ]; then

    docker build \
        --build-arg NOTARY_BRANCH=$NOTARY_BRANCH \
        --build-arg NOTARYPKG=$NOTARYPKG \
        --build-arg NOTARY_WRAPPER_BRANCH=$NOTARY_WRAPPER_BRANCH \
        --build-arg NOTARY_WRAPPER_PKG=$NOTARY_WRAPPER_PKG \
        -f ./$DOCKERFILE \
        -t $REGISTRY/notary-wrapper:$NOTARY_WRAPPER_TAG \
        --no-cache .
fi

# push notary-wrapper image
docker push $REGISTRY/notary-wrapper:$NOTARY_WRAPPER_TAG
