#!/usr/bin/env bash
# @Author zhendong.pan@nio.com
# GOARCH amd64
set -ex

cd $(dirname "$0")/..

TAG=v3.1.14
REGISTRY=adas-hub.nioint.com

rm -rf tmp && mkdir tmp

cd tmp && git clone https://git.nevint.com/adops/yearning-front && cd yearning-front && yarn && yarn build && cp -r dist ../../src/service/

cd ../..

docker build -t ${REGISTRY}/yearning:${TAG} -f Dockerfile .
if [[ ! $? -eq 0 ]]; then
    echo "yearning image built failed"
    exit 1
fi

docker push ${REGISTRY}/yearning:${TAG}


