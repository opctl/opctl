FROM docker:1.10.3

EXPOSE 42224

WORKDIR /opt/dev-op-spec

COPY engine ./
COPY tcp/swagger ./swagger/

RUN apk add -U py-pip && \
    pip install -U docker-compose==1.7.1 && \
    rm -rf /tmp/* /root/.cache

CMD ["/opt/dev-op-spec/engine"]
