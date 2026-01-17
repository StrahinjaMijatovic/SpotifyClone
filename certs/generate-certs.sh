#!/bin/bash
# Script to generate self-signed TLS certificates for development

echo "Generating self-signed TLS certificates..."

# Generate certificate valid for 365 days
openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
  -keyout key.pem \
  -out cert.pem \
  -subj "//CN=localhost" \
  -addext "subjectAltName=DNS:localhost,DNS:api-gateway,DNS:users-service,DNS:content-service,IP:127.0.0.1"

echo ""
echo "============================================"
echo "  TLS Certificates Generated Successfully!"
echo "============================================"
echo ""
echo "Files created:"
echo "  - cert.pem (public certificate)"
echo "  - key.pem  (private key - KEEP SECRET!)"
echo ""
echo "Valid for: 365 days"
echo "Valid for hosts: localhost, api-gateway, users-service, content-service"
echo ""
echo "To enable TLS in docker-compose.yml, set:"
echo "  TLS_ENABLED: \"true\""
echo ""
