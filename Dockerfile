FROM alpine:3.5

COPY /bin/slackwork /usr/bin/

ENTRYPOINT ["/usr/bin/slackwork"]
