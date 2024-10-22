#!/bin/bash
set -e

if [ ! -d /etc/dkimkeys ]; then
    mkdir -p /etc/dkimkeys || true
fi

if [[ -w /etc/dkimkeys ]]; then
    chmod -R 600 /etc/dkimkeys

    if [[ ! -f /etc/dkimkeys/mail.private ]]; then
        echo "Generating DKIM key..."
        opendkim-genkey -s mail -d "$DOMAIN" -D /etc/dkimkeys
        chown opendkim:opendkim /etc/dkimkeys/mail.private
    fi
fi

# Function to get external IP with retry
get_external_ip() {
    EXTERNAL_IP=$(curl -s --max-time 10 --retry 10 --retry-delay 5 https://api.ipify.org)
    if [ -n "$EXTERNAL_IP" ]; then
        return 0
    fi

    echo "Warning: Failed to retrieve external IP address after 10 attempts."
    return 1
}

# Get the external IP address with retry
if ! get_external_ip; then
    EXTERNAL_IP="YOUR_SERVER_IP"
    echo "Unable to automatically detect your external IP. Please replace 'YOUR_SERVER_IP' with your actual server IP."
fi

echo "================================================"
echo "SPF DNS record for $DOMAIN:"
echo "Add this to your DNS TXT record for: $DOMAIN"
echo ""
echo "v=spf1 ip4:$EXTERNAL_IP -all"
echo ""
echo "================================================"

echo ""
echo "Make sure to update your SPF record with the current external IP: $EXTERNAL_IP"
echo "This helps ensure your emails are not marked as spam."
echo "If 'YOUR_SERVER_IP' is shown, please replace it with your actual server IP."
echo ""

if [[ -f /etc/dkimkeys/mail.txt ]]; then
    echo "================================================"
    echo "DKIM DNS record for $DOMAIN:"
    # Extract the content inside the parentheses and trim whitespace
    echo "add this to your DNS TXT record for: mail._domainkey.$DOMAIN"
    echo ""
    cat /etc/dkimkeys/mail.txt
    echo ""
        echo "================================================"
fi

echo ""
echo "validating DKIM key, use following command:"
echo "opendkim-testkey -s mail -d $DOMAIN -v"
echo ""


echo "************************************************"

echo "================================================"
echo "Checking DNS records:"
echo ""

# Check MX record
echo "MX record for $DOMAIN:"
dig +short MX $DOMAIN
echo ""

# Check SPF record
echo "SPF record for $DOMAIN:"
dig +short TXT $DOMAIN
echo ""

# Check DKIM record
echo "DKIM record for mail._domainkey.$DOMAIN:"
dig +short TXT mail._domainkey.$DOMAIN
echo ""

echo "================================================"
echo ""
echo "If the records are not visible or incorrect, please ensure you have added"
echo "the SPF and DKIM records to your DNS as instructed above."
echo ""

echo "================================================"
echo "Testing DKIM key:"
opendkim-testkey -s mail -d $DOMAIN -v || true
echo "================================================"

[ ! -f /var/spool/postfix/etc/resolv.conf ] && cp /etc/resolv.conf /var/spool/postfix/etc/resolv.conf || true

/usr/sbin/opendkim
/usr/sbin/postfix start-fg
