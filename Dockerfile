FROM centos:7.4.1708

RUN mkdir -p /lingwei

WORKDIR /lingwei

ADD lingwei .
ADD config.ini config.ini

RUN chmod +x lingwei

ENTRYPOINT ./lingwei