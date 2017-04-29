FROM alpine
MAINTAINER jspc <james@zero-internet.org.uk>

EXPOSE 8000
ADD dogfucker-linux /dogfucker
ADD views /views

ENTRYPOINT ["/dogfucker"]
