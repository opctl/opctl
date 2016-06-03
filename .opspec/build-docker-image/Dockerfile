FROM docker/compose:1.7.1

EXPOSE 42224

WORKDIR /opt/opctl

COPY engine ./
COPY tcp/swagger ./swagger/

ENTRYPOINT ["sh","-c"]
CMD ["/opt/opctl/engine"]
