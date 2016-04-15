FROM docker:1.10.3

EXPOSE 8080

COPY dev-op-spec-engine /usr/local/bin
COPY swagger /usr/local/bin/swagger

RUN apk add -U py-pip && \
    pip install -U pip docker-compose==1.7.0 && \
    rm -rf /tmp/* /root/.cache

WORKDIR /usr/local/bin

CMD ["dev-op-spec-engine"]
