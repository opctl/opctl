FROM docker:1.13-dind

ARG CHANNEL

COPY entrypoint.sh /usr/local/bin/

RUN curl -L https://bin.equinox.io/c/4fmGop7rntx/opctl-beta-linux-amd64.tgz | tar -xzv -C /usr/local/bin
RUN opctl self-update -c "$CHANNEL"

ENTRYPOINT ["entrypoint.sh"]
