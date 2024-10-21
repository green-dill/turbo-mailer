#!/bin/bash
set -e

if [ ! -d /etc/dkimkeys ]; then
    mkdir -p /etc/dkimkeys
fi

if [ ! -f /etc/dkimkeys/mail.private ]; then
    opendkim-genkey -s mail -d $DOMAIN -D /etc/dkimkeys
    chown opendkim:opendkim /etc/dkimkeys/mail.private

    # test dkim key
    # opendkim-testkey -s mail -d $DOMAIN -vvv
fi

echo "================================================"
echo "DKIM DNS record for $DOMAIN:"
# Extract the content inside the parentheses and trim whitespace
echo "add this to your DNS TXT record for: mail._domainkey.$DOMAIN"
echo ""
cat /etc/dkimkeys/mail.txt
echo ""
echo "================================================"

echo ""
echo "validating DKIM key, use following command:"
echo "opendkim-testkey -s mail -d $DOMAIN -v"
echo ""


[ ! -f /var/spool/postfix/etc/resolv.conf ] && cp /etc/resolv.conf /var/spool/postfix/etc/resolv.conf || true

/usr/sbin/opendkim
/usr/sbin/postfix start-fg
