FROM alpine:3.10

RUN mkdir -p /lingwei

WORKDIR /lingwei

ADD lingwei .
ADD config.ini config.ini

RUN chmod +x lingwei

ENTRYPOINT ./lingwei
