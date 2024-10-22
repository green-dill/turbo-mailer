FROM debian:bookworm

RUN apt-get update && apt-get install -yq vim curl procps dnsutils postfix opendkim opendkim-tools rsyslog && \
    apt-get clean && apt-get autoremove && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p /etc/dkimkeys
RUN chown -R opendkim:opendkim /etc/dkimkeys && \
    chmod -R 600 /etc/dkimkeys


COPY config/header_checks /etc/postfix/header_checks
COPY config/aliases /etc/aliases
COPY config/rsyslog.conf /etc/rsyslog.conf
COPY start.sh /start.sh


CMD ["bash", "/start.sh"]
