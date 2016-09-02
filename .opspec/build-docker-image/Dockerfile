FROM docker/compose:1.8.0

EXPOSE 42224

WORKDIR /opt/opspec-io

COPY engine ./
COPY tcp/swagger ./swagger/

ENTRYPOINT ["sh","-c"]
CMD ["/opt/opspec-io/engine"]
