FROM docker:20.10.17-dind

COPY opctl /usr/local/bin/
COPY entrypoint.sh /usr/local/bin/

ENTRYPOINT ["entrypoint.sh"]
