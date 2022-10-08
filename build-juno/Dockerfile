FROM alpine:3.15
LABEL maintainer="panzhendong <zhendong.pan@nio.com>"

COPY Juno /opt/Juno
COPY conf.toml /opt/conf.toml

WORKDIR /opt/

EXPOSE 50001

CMD ["/opt/Juno"]
