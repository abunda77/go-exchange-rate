# Build Information - Exchange Rate App

Generated: 2026-01-27

## Available Versions

### 1. Standard Version
- **File**: `exchange-rate-app.exe`
- **Size**: ~10.4 MB
- **Requires**: `.env` file with API_KEY
- **Port**: 8889
- **Source**: `main.go` + `api/client.go`
- **Security**: ✅ API key terpisah dari binary

### 2. Embedded Version ⭐ NEW
- **File**: `exchange-rate-app-embedded.exe`
- **Size**: ~10.0 MB
- **Requires**: Nothing (standalone)
- **Port**: 8890
- **Source**: `main_embedded.go` + `api/client_embedded.go`
- **Security**: ⚠️ API key hardcoded dalam binary

## Quick Start

### Standard Version:
```powershell
# Pastikan ada file .env
.\exchange-rate-app.exe
```

### Embedded Version:
```powershell
# Tidak perlu .env file
.\exchange-rate-app-embedded.exe
```

## Rebuild

### Standard:
```powershell
go build -o exchange-rate-app.exe main.go
# atau
.\build.ps1
```

### Embedded:
```powershell
go build -ldflags="-H windowsgui" -o exchange-rate-app-embedded.exe main_embedded.go
# atau
.\build-embedded.ps1
```

## Files Created

```
✅ api/client_embedded.go       - API client dengan embedded key
✅ main_embedded.go             - Main app versi embedded
✅ build-embedded.ps1           - Build script
✅ VERSION_COMPARISON.md        - Dokumentasi perbandingan
✅ exchange-rate-app-embedded.exe - Compiled binary
✅ README.md (updated)          - Updated dengan info embedded version
```

## Distribution

### For End Users (Non-Technical):
Use **Embedded Version** - Cukup distribute file exe saja.

### For Developers/Technical Users:
Use **Standard Version** - Provide exe + .env.example, user set their own API key.

## Important Notes

⚠️ **Security Warning**: Versi embedded memiliki API key di-hardcode. Jangan distribute secara publik jika API key sensitive.

✅ **Best Practice**: Untuk production, gunakan standard version dengan proper secret management.

📝 **Documentation**: Lihat `VERSION_COMPARISON.md` untuk detail lengkap perbedaan kedua versi.

---

Built with Go + WebView2 ❤️
