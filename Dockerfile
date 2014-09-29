############################################################
# Dockerfile to run Neverdown inside a container
# Based on Ubuntu Image
############################################################
FROM ubuntu
MAINTAINER Thomas Sileo
EXPOSE 3000
WORKDIR /data/useragentparser
RUN apt-get update
RUN apt-get -y install wget
RUN wget https://raw.githubusercontent.com/tobie/ua-parser/master/regexes.yaml
ADD ./useragentparser /opt/useragentparser
ENTRYPOINT /opt/useragentparser
