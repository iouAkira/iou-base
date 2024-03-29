FROM alpine as goBuild
LABEL AUTHOR="iouAkira <ZS5ha2ltb3RvLmFraXJhQGdtYWlsLmNvbQ==>"
RUN set -ex \
    && apk update \
    && apk upgrade \
    && apk add --no-cache go git

RUN set -ex \
    && cd / \
    && git clone https://github.com/iouAkira/iou-base.git \
    && cd /iou-base/repo_sync \
    && go build repo_sync.go

FROM alpine
LABEL AUTHOR="iouAkira <ZS5ha2ltb3RvLmFraXJhQGdtYWlsLmNvbQ==>"

ENV VER=0.1 \
    MNT_DIR=/data \
    REPOS_DIR=/iouRepos \
    LOCAL_DIR=/iouLocalDir \
    CRON_FILE_DIR=/iouCron

RUN set -ex \
    && apk update && apk upgrade\
    && apk add --no-cache tzdata git jq\
    && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && mkdir /data \
    && mkdir /iouCron \
    && mkdir /iouRepos

COPY --from=goBuild /iou-base/repo_sync/repo_sync /usr/local/bin/repo_sync
COPY entrypoint.sh /usr/local/bin/entrypoint.sh

RUN chmod +x /usr/local/bin/repo_sync \ 
    && chmod +x /usr/local/bin/entrypoint.sh

VOLUME [ "/data" ]

WORKDIR /iouRepos

ENTRYPOINT ["entrypoint.sh"]

CMD ["keep-run"]

