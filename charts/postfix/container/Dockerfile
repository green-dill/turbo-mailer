FROM debian:bookworm

RUN apt-get update && apt-get install -yq vim curl procps dnsutils postfix opendkim opendkim-tools

COPY config/header_checks /etc/postfix/header_checks
COPY config/aliases /etc/aliases

# CMD ["/usr/sbin/postfix", "start-fg"]

CMD \
  [ ! -f /var/spool/postfix/etc/resolv.conf ] && cp /etc/resolv.conf /var/spool/postfix/etc/resolv.conf && \
  /usr/sbin/postfix start-fg
