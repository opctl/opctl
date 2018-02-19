FROM docker:18.02.0-dind

ARG CHANNEL

COPY entrypoint.sh /usr/local/bin/

ADD https://bin.equinox.io/c/4fmGop7rntx/opctl-${CHANNEL}-linux-amd64.tgz /tmp
RUN tar -C /usr/local/bin -xvzf /tmp/opctl-${CHANNEL}-linux-amd64.tgz opctl
RUN opctl self-update -c "$CHANNEL"

ENTRYPOINT ["entrypoint.sh"]
