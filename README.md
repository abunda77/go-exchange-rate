# Exchange Rate GUI Application

Aplikasi desktop berbasis Go untuk menampilkan real-time exchange rates menggunakan API dari use.api.co.id.

## Features

- ✨ Modern GUI dengan dark theme dan glass-morphism effects
- 💱 Menampilkan exchange rates untuk 19 currencies (AED, AUD, CAD, CHF, CNY, EUR, GBP, HKD, IDR, INR, JPY, KRW, MYR, NZD, PHP, SAR, SGD, THB, USD)
- 🔍 3 mode viewing:
  - **All Rates** - Lihat semua exchange rates tersedia
  - **Currency Lookup** - Cari rate berdasarkan base currency
  - **Pair Exchange** - Cari rate untuk currency pair tertentu
- ⚡ Fast & lightweight - No CGO required
- 🔄 Auto-refresh data dari API

## Tech Stack

- **Language**: Go 1.20+
- **GUI Framework**: [jchv/go-webview2](https://github.com/jchv/go-webview2)
- **Runtime**: Microsoft Edge WebView2
- **API**: [use.api.co.id](https://use.api.co.id)

## Prerequisites

- **Windows 10/11** dengan Microsoft Edge WebView2 Runtime
  - WebView2 Runtime biasanya sudah terinstall di Windows 10/11 yang up-to-date
  - Jika belum ada, download dari [Microsoft](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)
- **API Key** dari use.api.co.id

## Installation

### Download Binary (Recommended)

1. Download `exchange-rate-app.exe` dari releases
2. Copy file `.env.example` menjadi `.env`
3. Edit `.env` dan masukkan API key Anda
4. Jalankan `exchange-rate-app.exe`

### Build from Source

```bash
# Clone repository
git clone <repository-url>
cd window-app

# Copy environment file
cp .env.example .env

# Edit .env dan masukkan API key
# API_KEY=your_api_key_here

# Build
go build -o exchange-rate-app.exe

# Run
.\exchange-rate-app.exe
```

## Configuration

Buat file `.env` di root directory dengan isi:

```env
API_KEY=your_api_key_here
```

Replace `your_api_key_here` dengan API key yang valid dari use.api.co.id.

## Usage

### All Rates Tab
1. Klik tombol **"Refresh Rates"**
2. Aplikasi akan menampilkan semua exchange rates untuk 19 currencies
3. Data ditampilkan dalam format: Base currency → Target currencies

### Currency Lookup Tab
1. Input kode currency (3 karakter) - contoh: `USD`, `EUR`, `IDR`
2. Klik **"Get Rate"**
3. Aplikasi menampilkan exchange rates untuk currency tersebut

### Pair Exchange Tab
1. Input currency pair - contoh: `USDIDR`
2. Untuk multiple pairs, gunakan comma - contoh: `USDIDR,EURIDR,MYRIDR`
3. Klik **"Get Rate"**
4. Aplikasi menampilkan rate untuk pair yang diminta

## API Endpoints

Aplikasi menggunakan 3 endpoints dari use.api.co.id:

- `GET /api/exchange-rates` - Semua exchange rates
- `GET /currency/:currency` - Rate untuk base currency tertentu
- `GET /currency/exchange-rate?pair=PAIR` - Rate untuk currency pair

## Project Structure

```
window-app/
├── api/
│   └── client.go          # HTTP client untuk API
├── main.go                # Main application + GUI
├── .env                   # API key configuration (gitignored)
├── .env.example           # Template untuk .env
├── go.mod                 # Go module dependencies
├── go.sum                 # Dependencies checksum
├── planning.md            # Project planning document
└── README.md              # This file
```

## Troubleshooting

### Error: Invalid API key (401)
- Pastikan file `.env` berisi API key yang valid
- Restart aplikasi setelah mengubah `.env`

### Aplikasi menampilkan layar kosong/blank
- Pastikan Microsoft Edge WebView2 Runtime terinstall
- Update Windows ke versi terbaru
- Download WebView2 Runtime dari [Microsoft](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)

### Port 8889 already in use
- Tutup instance aplikasi yang lain
- Atau ubah port di `main.go` (line ~66)

## Development

### Build for Production

```bash
go build -ldflags="-s -w" -o exchange-rate-app.exe
```

### Run Tests

```bash
go test ./...
```

## Dependencies

- `github.com/jchv/go-webview2` - WebView2 bindings untuk Go
- `github.com/joho/godotenv` - Load .env files

## License

This project is open source and available under the MIT License.

## Author

Built with ❤️ using Go and WebView2

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
