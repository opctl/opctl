FROM --platform=$TARGETPLATFORM ghcr.io/linuxcontainers/alpine:3.11.3
ARG TARGETOS TARGETARCH
COPY opctl-$TARGETOS-$TARGETARCH /usr/local/bin/opctl
CMD [ "opctl", "node", "create" ]
