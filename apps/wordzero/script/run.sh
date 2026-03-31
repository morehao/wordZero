#!/bin/bash

WORK_ENV=prod
APP=wordzero
DOMAIN=yygu.cn
DATA_DIR=/data/${APP}.${WORK_ENV}/

IMAGE_NAME=registry.cn-beijing.aliyuncs.com/wa/openrpa
IMAGE_TAG=${APP}_v0.8.8-130
IMAGE=${IMAGE_NAME}:${IMAGE_TAG}

set -ex;

docker pull ${IMAGE};

docker service rm ${APP}-${WORK_ENV} || true
docker service create \
    --name=${APP}-${WORK_ENV} \
    --network=bridge \
    --replicas=1 \
    --replicas-max-per-node=1 \
    --mount=type=bind,source=${DATA_DIR},target=${DATA_DIR} \
    --container-label=traefik.enable=true \
    --container-label=traefik.http.routers.rt-https-${APP}-${WORK_ENV}.entrypoints=websecure \
    --container-label=traefik.http.routers.rt-https-${APP}-${WORK_ENV}.rule=Host\(\`${DOMAIN}\`\)\&\&PathPrefix\(\`/v1/wordzero\`\) \
    --container-label=traefik.http.routers.rt-https-${APP}-${WORK_ENV}.tls.certResolver=letsencrypt \
    --container-label=traefik.http.routers.rt-https-${APP}-${WORK_ENV}.tls.domains[0].main=${DOMAIN} \
    --container-label=traefik.http.services.svc-${APP}-${WORK_ENV}.loadbalancer.server.port=8080 \
    ${IMAGE} -c ${DATA_DIR}/conf/config.yaml