FROM alpine
MAINTAINER jspc <james@zero-internet.org.uk>

EXPOSE 8000
ADD dogf-linux /dogf
ADD views /views

ENTRYPOINT ["/dogf"]
