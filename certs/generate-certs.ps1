# PowerShell script to generate self-signed TLS certificates for development
# Run: .\generate-certs.ps1

Write-Host "Generating self-signed TLS certificates..." -ForegroundColor Cyan

# Generate certificate valid for 365 days
openssl req -x509 -nodes -days 365 -newkey rsa:2048 `
  -keyout key.pem `
  -out cert.pem `
  -subj "//CN=localhost"

Write-Host ""
Write-Host "============================================" -ForegroundColor Green
Write-Host "  TLS Certificates Generated Successfully!" -ForegroundColor Green
Write-Host "============================================" -ForegroundColor Green
Write-Host ""
Write-Host "Files created:"
Write-Host "  - cert.pem (public certificate)" -ForegroundColor Yellow
Write-Host "  - key.pem  (private key - KEEP SECRET!)" -ForegroundColor Red
Write-Host ""
Write-Host "Valid for: 365 days"
Write-Host ""
Write-Host "To enable TLS, run:" -ForegroundColor Cyan
Write-Host "  docker-compose -f docker-compose.yml -f docker-compose.tls.yml up -d"
Write-Host ""
