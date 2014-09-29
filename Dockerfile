############################################################
# Dockerfile to run Neverdown inside a container
# Based on Ubuntu Image
############################################################
FROM ubuntu
MAINTAINER Thomas Sileo
EXPOSE 3000
WORKDIR /data/useragentparser
ADD ./regexes.yaml /data/useragentparser/regexes.yaml
ADD ./useragentparser /opt/useragentparser
ENTRYPOINT /opt/useragentparser
