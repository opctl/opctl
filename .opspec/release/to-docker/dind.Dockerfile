FROM docker:19.03-dind

COPY opctl /usr/local/bin/
COPY entrypoint.sh /usr/local/bin/

ENTRYPOINT ["entrypoint.sh"]
