#!/bin/bash

# Script to generate self-signed TLS certificates for development/testing
# For production, use certificates from a trusted CA (e.g., Let's Encrypt)

set -e

CERT_DIR="certs"
DAYS_VALID=365
KEY_SIZE=4096

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}=== TLS Certificate Generator ===${NC}"
echo ""

# Create certs directory if it doesn't exist
mkdir -p $CERT_DIR

# Check if OpenSSL is installed
if ! command -v openssl &> /dev/null; then
    echo -e "${RED}Error: OpenSSL is not installed${NC}"
    echo "Please install OpenSSL first:"
    echo "  Ubuntu/Debian: sudo apt-get install openssl"
    echo "  CentOS/RHEL:   sudo yum install openssl"
    echo "  macOS:         brew install openssl"
    echo "  Windows:       Download from https://slproweb.com/products/Win32OpenSSL.html"
    exit 1
fi

echo -e "${YELLOW}Generating RSA private key (${KEY_SIZE} bits)...${NC}"
openssl genrsa -out $CERT_DIR/key.pem $KEY_SIZE

echo -e "${YELLOW}Generating self-signed certificate (valid for ${DAYS_VALID} days)...${NC}"
openssl req -new -x509 \
    -key $CERT_DIR/key.pem \
    -out $CERT_DIR/cert.pem \
    -days $DAYS_VALID \
    -subj "/C=RS/ST=Serbia/L=Belgrade/O=SpotifyClone/OU=Development/CN=localhost" \
    -addext "subjectAltName=DNS:localhost,DNS:api-gateway,DNS:users-service,DNS:content-service,DNS:ratings-service,DNS:subscriptions-service,DNS:notifications-service,DNS:recommendation-service,IP:127.0.0.1"

# Set proper permissions
chmod 600 $CERT_DIR/key.pem
chmod 644 $CERT_DIR/cert.pem

echo ""
echo -e "${GREEN}=== Certificates generated successfully! ===${NC}"
echo ""
echo "Files created:"
echo "  - $CERT_DIR/cert.pem (certificate)"
echo "  - $CERT_DIR/key.pem (private key)"
echo ""
echo "Certificate details:"
openssl x509 -in $CERT_DIR/cert.pem -noout -subject -dates
echo ""
echo -e "${YELLOW}NOTE: These are self-signed certificates for development only!${NC}"
echo "For production, use certificates from a trusted CA."
echo ""

# Copy certificates to service directories
echo -e "${YELLOW}Copying certificates to service directories...${NC}"

# List of services that need certificates
SERVICES=("api-gateway" "users-service" "content-service" "ratings-service" "subscriptions-service" "notifications-service" "recommendation-service")

for service in "${SERVICES[@]}"; do
    if [ -d "$service" ]; then
        mkdir -p "$service/certs"
        cp $CERT_DIR/cert.pem "$service/certs/"
        cp $CERT_DIR/key.pem "$service/certs/"
        echo "  - Copied to $service/certs/"
    fi
done

echo ""
echo -e "${GREEN}Done!${NC}"
echo ""
echo "To enable HTTPS, set the following environment variables:"
echo "  TLS_ENABLED=true"
echo "  TLS_CERT_FILE=certs/cert.pem"
echo "  TLS_KEY_FILE=certs/key.pem"
