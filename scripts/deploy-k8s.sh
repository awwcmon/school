helm repo add ${CHART_REPO_NAME} ${CHART_URL}
helm repo update
helm upgrade \
--install ${params.IMAGE_NAME}${params.RELEASE_NAME} ${CHART_REPO_NAME}/${params.IMAGE_NAME} \
--namespace ${env.NAMESPACE}