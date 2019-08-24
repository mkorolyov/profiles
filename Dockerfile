FROM alpine:latest

COPY bin/linux_profiles /profiles
RUN chmod u+x /profiles
COPY etc/config.yaml /config.yaml

ENTRYPOINT ["/profiles"]