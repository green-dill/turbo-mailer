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

# Function to get external IP with retry
get_external_ip() {
    local max_attempts=10
    local attempt=1
    local delay=5

    while [ $attempt -le $max_attempts ]; do
        EXTERNAL_IP=$(curl -s --max-time 10 https://api.ipify.org)
        if [ -n "$EXTERNAL_IP" ]; then
            return 0
        fi
        echo "Attempt $attempt failed. Retrying in $delay seconds..."
        sleep $delay
        attempt=$((attempt + 1))
    done

    echo "Warning: Failed to retrieve external IP address after $max_attempts attempts."
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
