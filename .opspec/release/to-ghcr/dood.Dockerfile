FROM --platform=$TARGETPLATFORM ghcr.io/linuxcontainers/alpine:3.18
ARG TARGETOS TARGETARCH
COPY opctl-$TARGETOS-$TARGETARCH /usr/local/bin/opctl
EXPOSE 42224/tcp
CMD [ "opctl", "node", "create" ]
