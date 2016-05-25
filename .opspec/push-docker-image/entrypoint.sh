#!/bin/sh
docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD} -e . && \
docker push ${DOCKER_REPO_NAME}
