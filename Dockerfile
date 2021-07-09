FROM alpine
LABEL AUTHOR="iouAkira <ZS5ha2ltb3RvLmFraXJhQGdtYWlsLmNvbQ==>"

ENV VER=0.1 \
    MNT_DIR=/data \
    CRON_FILE_DIR=/iouCron \
    REPOS_DIR=/iouRepos

RUN set -ex \
    && apk update && apk upgrade\
    && apk add --no-cache tzdata git jq\
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && mkdir /data \
    && mkdir /iouCron \
    && mkdir /iouRepos

COPY entrypoint.sh /usr/local/bin/docker_entrypoint.sh

RUN chmod +x /usr/local/bin/docker_entrypoint.sh

WORKDIR /iouRepos

ENTRYPOINT ["docker_entrypoint.sh","keep-run"]

