#!/usr/bin/env bash

NAMESPACES='
default
dev
knative-build
'

CLUSTER='docker-for-desktop-cluster'
K8S_USER='docker-for-desktop'

for NAMESPACE in ${NAMESPACES}; do
  CONTEXT=${NAMESPACE}-${CLUSTER}
  kubectl config set-context ${CONTEXT} \
    --cluster ${CLUSTER} \
    --user ${K8S_USER} \
    --namespace ${NAMESPACE}
done

