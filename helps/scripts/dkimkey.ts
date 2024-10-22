import * as crypto from 'crypto';
import * as readline from 'readline';

function generateDKIMKey(domain: string) {
    // Generate an RSA key pair
    const { privateKey, publicKey } = crypto.generateKeyPairSync('rsa', {
        modulusLength: 2048, // Key size
        publicKeyEncoding: {
            type: 'spki',
            format: 'pem'
        },
        privateKeyEncoding: {
            type: 'pkcs8',
            format: 'pem'
        }
    });

    // Generate DKIM DNS record with h=sha256
    const dkimRecord = `v=DKIM1; k=rsa; h=sha256; p=${publicKey.toString().replace(/-----BEGIN PUBLIC KEY-----|-----END PUBLIC KEY-----|\n/g, '')}`;

    console.log(`\nAdd TXT record to: mail._domainkey.${domain}\n\n`);
    console.log(`${dkimRecord}\n\n`);
    console.log(`Generated Private Key:\n\n${privateKey}`);
}

const rl = readline.createInterface({
    input: process.stdin,
    output: process.stdout
});

// Prompt user for domain input
rl.question('Please enter your domain (default: mail.example.com): ', (input) => {
    const domain = input.trim() || 'mail.example.com';
    generateDKIMKey(domain);
    rl.close();
});
