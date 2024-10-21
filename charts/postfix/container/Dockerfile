FROM debian:bookworm

RUN apt-get update && apt-get install -yq vim curl procps dnsutils postfix opendkim opendkim-tools python3 python3-pip && \
    pip3 install dnspython requests && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/dkimkeys
RUN chown -R opendkim:opendkim /etc/dkimkeys && \
    chmod -R 600 /etc/dkimkeys


COPY config/header_checks /etc/postfix/header_checks
COPY config/aliases /etc/aliases
COPY start.sh /start.sh


CMD ["bash", "/start.sh"]
