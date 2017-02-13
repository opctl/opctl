FROM alpine:3.5

ARG VERSION

EXPOSE 42224

RUN apk add --no-cache curl
RUN curl -L https://bin.equinox.io/a/5Cw5dTJ3ECp/opctl-$VERSION-linux-amd64.tar.gz | tar -xzv -C /usr/local/bin

ENTRYPOINT ["opctl", "node", "create"]
