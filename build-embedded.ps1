# Build Embedded Version (No Environment Required)

Write-Host "Building embedded version of Exchange Rate App..." -ForegroundColor Cyan

# Build the embedded version using build tags
Write-Host "`nCompiling with embedded tag..." -ForegroundColor Yellow
go build -tags embedded -ldflags="-H windowsgui" -o exchange-rate-app-embedded.exe

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ Build successful!" -ForegroundColor Green
    Write-Host "`nGenerated file:" -ForegroundColor Cyan
    Write-Host "  - exchange-rate-app-embedded.exe (with embedded API key)" -ForegroundColor White
    Write-Host "`nThis version does NOT require .env file" -ForegroundColor Yellow
    Write-Host "API key is hardcoded in the binary." -ForegroundColor Yellow
} else {
    Write-Host "`n❌ Build failed!" -ForegroundColor Red
    exit 1
}
