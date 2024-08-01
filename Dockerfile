FROM vechain/thor:v2.1.3

USER root

RUN apk update && apk upgrade && apk add curl

ENTRYPOINT ["/bin/sh", "-c", "thor"]
