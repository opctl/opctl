#!/bin/sh
docker login -u ${DOCKER_USERNAME} -p ${DOCKER_PASSWORD} -e ${DOCKER_EMAIL} && \
docker push ${DOCKER_REPO_NAME}:$(cat target/VERSION)