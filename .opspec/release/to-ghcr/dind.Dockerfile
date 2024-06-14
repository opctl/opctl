FROM --platform=$TARGETPLATFORM docker:20.10.17-dind
ARG TARGETOS TARGETARCH
COPY opctl-$TARGETOS-$TARGETARCH /usr/local/bin/opctl
COPY entrypoint.sh /usr/local/bin/

ENTRYPOINT ["entrypoint.sh"]
