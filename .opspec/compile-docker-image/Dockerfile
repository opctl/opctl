FROM docker:17.07.0-dind

ARG CHANNEL

COPY entrypoint.sh /usr/local/bin/

ADD https://bin.equinox.io/c/4fmGop7rntx/opctl-beta-linux-amd64.tgz /tmp
RUN mv /tmp/opctl /usr/local/bin && opctl self-update -c "$CHANNEL"

ENTRYPOINT ["entrypoint.sh"]
