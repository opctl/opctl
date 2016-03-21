FROM docker:1.10.3

EXPOSE 8080

COPY dev-op-spec-engine /usr/local/bin

RUN apk add -U git py-pip && \
    pip install -U pip docker-compose==1.6.0 && \
    rm -rf /tmp/* /root/.cache

CMD ["dev-op-spec-engine"]