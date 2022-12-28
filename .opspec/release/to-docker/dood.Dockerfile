FROM --platform=$TARGETPLATFORM alpine:3.11.3
ARG TARGETOS TARGETARCH
COPY opctl-$TARGETOS-$TARGETARCH /usr/local/bin/opctl
EXPOSE 42224/tcp
CMD [ "opctl", "node", "create" ]
