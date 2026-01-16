# PowerShell script to generate self-signed TLS certificates for development/testing
# For production, use certificates from a trusted CA (e.g., Let's Encrypt)

$ErrorActionPreference = "Stop"

$CERT_DIR = "certs"
$DAYS_VALID = 365
$KEY_SIZE = 4096

Write-Host "=== TLS Certificate Generator ===" -ForegroundColor Green
Write-Host ""

# Create certs directory if it doesn't exist
if (-not (Test-Path $CERT_DIR)) {
    New-Item -ItemType Directory -Path $CERT_DIR | Out-Null
}

# Check if OpenSSL is available
$opensslPath = $null
$possiblePaths = @(
    "openssl",
    "C:\Program Files\OpenSSL-Win64\bin\openssl.exe",
    "C:\Program Files (x86)\OpenSSL-Win32\bin\openssl.exe",
    "C:\OpenSSL-Win64\bin\openssl.exe"
)

foreach ($path in $possiblePaths) {
    try {
        $null = & $path version 2>$null
        $opensslPath = $path
        break
    } catch {
        continue
    }
}

if (-not $opensslPath) {
    Write-Host "Error: OpenSSL is not installed or not in PATH" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install OpenSSL:"
    Write-Host "  1. Download from: https://slproweb.com/products/Win32OpenSSL.html"
    Write-Host "  2. Or use chocolatey: choco install openssl"
    Write-Host "  3. Or use winget: winget install OpenSSL"
    exit 1
}

Write-Host "Using OpenSSL: $opensslPath" -ForegroundColor Cyan
Write-Host ""

# Create OpenSSL config file for SAN (Subject Alternative Names)
$opensslConfig = @"
[req]
default_bits = $KEY_SIZE
prompt = no
default_md = sha256
distinguished_name = dn
req_extensions = req_ext
x509_extensions = v3_ca

[dn]
C = RS
ST = Serbia
L = Belgrade
O = SpotifyClone
OU = Development
CN = localhost

[req_ext]
subjectAltName = @alt_names

[v3_ca]
subjectAltName = @alt_names
basicConstraints = critical, CA:true

[alt_names]
DNS.1 = localhost
DNS.2 = api-gateway
DNS.3 = users-service
DNS.4 = content-service
DNS.5 = ratings-service
DNS.6 = subscriptions-service
DNS.7 = notifications-service
DNS.8 = recommendation-service
IP.1 = 127.0.0.1
"@

$configPath = "$CERT_DIR\openssl.cnf"
$opensslConfig | Out-File -FilePath $configPath -Encoding ASCII

Write-Host "Generating RSA private key ($KEY_SIZE bits)..." -ForegroundColor Yellow
& $opensslPath genrsa -out "$CERT_DIR\key.pem" $KEY_SIZE 2>$null

Write-Host "Generating self-signed certificate (valid for $DAYS_VALID days)..." -ForegroundColor Yellow
& $opensslPath req -new -x509 `
    -key "$CERT_DIR\key.pem" `
    -out "$CERT_DIR\cert.pem" `
    -days $DAYS_VALID `
    -config $configPath 2>$null

# Clean up config file
Remove-Item $configPath -ErrorAction SilentlyContinue

Write-Host ""
Write-Host "=== Certificates generated successfully! ===" -ForegroundColor Green
Write-Host ""
Write-Host "Files created:"
Write-Host "  - $CERT_DIR\cert.pem (certificate)"
Write-Host "  - $CERT_DIR\key.pem (private key)"
Write-Host ""

# Show certificate details
Write-Host "Certificate details:" -ForegroundColor Cyan
& $opensslPath x509 -in "$CERT_DIR\cert.pem" -noout -subject -dates

Write-Host ""
Write-Host "NOTE: These are self-signed certificates for development only!" -ForegroundColor Yellow
Write-Host "For production, use certificates from a trusted CA."
Write-Host ""

# Copy certificates to service directories
Write-Host "Copying certificates to service directories..." -ForegroundColor Yellow

$services = @(
    "api-gateway",
    "users-service",
    "content-service",
    "ratings-service",
    "subscriptions-service",
    "notifications-service",
    "recommendation-service"
)

foreach ($service in $services) {
    if (Test-Path $service) {
        $serviceCertDir = "$service\certs"
        if (-not (Test-Path $serviceCertDir)) {
            New-Item -ItemType Directory -Path $serviceCertDir | Out-Null
        }
        Copy-Item "$CERT_DIR\cert.pem" "$serviceCertDir\" -Force
        Copy-Item "$CERT_DIR\key.pem" "$serviceCertDir\" -Force
        Write-Host "  - Copied to $serviceCertDir\"
    }
}

Write-Host ""
Write-Host "Done!" -ForegroundColor Green
Write-Host ""
Write-Host "To enable HTTPS, set the following environment variables:"
Write-Host "  TLS_ENABLED=true"
Write-Host "  TLS_CERT_FILE=certs/cert.pem"
Write-Host "  TLS_KEY_FILE=certs/key.pem"
