#!/bin/bash
CHART_REPO_NAME=$1
CHART_URL=$2
IMAGE_NAME=$3
NAMESPACE=$4
RELEASE_NAME=$5
#export KUBECONFIG=/home/jenkins/.kube/kubeconfig.yaml
helm repo add ${CHART_REPO_NAME} ${CHART_URL}
helm repo update
helm upgrade \
--install ${IMAGE_NAME}${RELEASE_NAME}  ${CHART_REPO_NAME}/${IMAGE_NAME} \
--namespace=${NAMESPACE}
