# Perbandingan Versi Standard vs Embedded

## Versi Standard (`exchange-rate-app.exe`)

### Kelebihan:
- ✅ **Lebih Aman** - API key disimpan di file `.env` terpisah, tidak di-hardcode dalam binary
- ✅ **Mudah Update** - Bisa mengganti API key tanpa perlu rebuild aplikasi
- ✅ **Best Practice** - Mengikuti 12-factor app methodology untuk konfigurasi

### Kekurangan:
- ❌ **Memerlukan .env** - Harus ada file `.env` di folder yang sama dengan exe
- ❌ **Distribusi Kompleks** - Harus mendistribusikan 2 file (exe + .env)

### Kapan Menggunakan:
- Development dan testing
- Deployment di environment yang bisa manage secrets
- Ketika perlu flexibility untuk ganti API key

### Build Command:
```powershell
.\build.ps1
```

### Run:
```powershell
.\exchange-rate-app.exe
```

**Port:** 8889

---

## Versi Embedded (`exchange-rate-app-embedded.exe`)

### Kelebihan:
- ✅ **Standalone** - Tidak memerlukan file tambahan (.env tidak diperlukan)
- ✅ **Mudah Distribusi** - Cukup distribusikan 1 file exe saja
- ✅ **User Friendly** - User tidak perlu konfigurasi apapun

### Kekurangan:
- ❌ **Kurang Aman** - API key di-hardcode dalam binary (bisa di-extract)
- ❌ **Susah Update** - Untuk ganti API key harus rebuild dan redistribute
- ❌ **Security Risk** - Jika file exe bocor, API key juga bocor

### Kapan Menggunakan:
- Distribusi ke end-user yang tidak technical
- Quick demo atau POC (Proof of Concept)
- Internal tools dengan API key yang tidak sensitive
- Environment dimana manage .env file sulit

### Build Command:
```powershell
.\build-embedded.ps1
```

### Run:
```powershell
.\exchange-rate-app-embedded.exe
```

**Port:** 8890

---

## Perbedaan Teknis

| Aspek | Standard | Embedded |
|-------|----------|----------|
| File Source | `main.go` | `main_embedded.go` |
| API Client | `api.Client` | `api.ClientEmbedded` |
| Environment | `godotenv` + `os.Getenv()` | Hardcoded constant |
| Port | 8889 | 8890 |
| Dependencies | `godotenv` required | No `godotenv` needed |
| Title | "Exchange Rate App" | "Exchange Rate App (Embedded API)" |
| Badge | None | "EMBEDDED API" |

---

## Rekomendasi Keamanan

### ⚠️ PENTING untuk Versi Embedded:

1. **Jangan Commit** - Jangan commit `main_embedded.go` dengan API key asli ke Git
2. **Rotate API Key** - Jika exe tersebar luas, segera rotate API key
3. **Use Obfuscation** - Pertimbangkan untuk obfuscate binary dengan tools seperti:
   - `garble` - https://github.com/burrowers/garble
   - `go build` dengan flags khusus
4. **Limit Distribution** - Hanya distribute ke pihak yang dipercaya
5. **Monitor Usage** - Monitor penggunaan API key untuk detect abuse

### Alternative Security Approach:

Jika memerlukan standalone app tapi tetap aman:
- Gunakan API key dengan permission terbatas
- Implement rate limiting di API side
- Gunakan short-lived tokens instead of long-lived API keys
- Build licensing/activation system

---

## Source Code Location

```
window-app/
├── api/
│   ├── client.go              # Client untuk versi standard
│   └── client_embedded.go     # Client untuk versi embedded
├── main.go                    # Entry point versi standard
├── main_embedded.go           # Entry point versi embedded
├── build.ps1                  # Build script untuk standard
└── build-embedded.ps1         # Build script untuk embedded
```

---

## Build Both Versions

Untuk build kedua versi sekaligus:

```powershell
# Build standard version
.\build.ps1

# Build embedded version
.\build-embedded.ps1

# Hasil:
# - exchange-rate-app.exe (standard)
# - exchange-rate-app-embedded.exe (embedded)
```

---

## Testing

Untuk test kedua versi bisa dijalankan bersamaan karena menggunakan port berbeda:

```powershell
# Terminal 1 - Standard version
.\exchange-rate-app.exe

# Terminal 2 - Embedded version (bisa berjalan bersamaan)
.\exchange-rate-app-embedded.exe
```

Standard version akan buka di port 8889, Embedded di port 8890.
