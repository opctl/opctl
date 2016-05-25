FROM docker:1.11.1

EXPOSE 42224

WORKDIR /opt/opctl

COPY engine ./
COPY tcp/swagger ./swagger/

RUN apk add -U py-pip && \
    pip install -U docker-compose==1.7.1 && \
    rm -rf /tmp/* /root/.cache

CMD ["/opt/opctl/engine"]
