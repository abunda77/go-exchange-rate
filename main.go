package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jchv/go-webview2"
	"github.com/joho/godotenv"

	"window-app/api"
)

var apiClient *api.Client

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Create API client
	apiClient = api.NewClient()

	// Start local HTTP server first
	go startHTTPServer()

	// Wait for server to start
	time.Sleep(500 * time.Millisecond)

	// Create webview
	w := webview2.NewWithOptions(webview2.WebViewOptions{
		Debug:     true,
		AutoFocus: true,
		WindowOptions: webview2.WindowOptions{
			Title:  "Exchange Rate App",
			Width:  900,
			Height: 650,
			IconId: 0,
			Center: true,
		},
	})
	if w == nil {
		log.Fatal("Failed to create webview. Make sure Microsoft Edge WebView2 Runtime is installed.")
	}
	defer w.Destroy()

	// Bind Go functions to JavaScript
	w.Bind("getAllRates", getAllRates)
	w.Bind("getCurrencyRate", getCurrencyRate)
	w.Bind("getPairRate", getPairRate)

	// Navigate to local server
	w.Navigate("http://localhost:8889/")
	w.Run()
}

func startHTTPServer() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/api/all-rates", handleAllRates)
	http.HandleFunc("/api/currency/", handleCurrencyRate)
	http.HandleFunc("/api/pair", handlePairRate)
	log.Println("Starting HTTP server on :8889")
	if err := http.ListenAndServe(":8889", nil); err != nil {
		log.Printf("HTTP server error: %v", err)
	}
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, getHTMLContent())
}

func handleAllRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	response, err := apiClient.GetAllExchangeRates()
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(response)
}

func handleCurrencyRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	currency := strings.TrimPrefix(r.URL.Path, "/api/currency/")
	response, err := apiClient.GetCurrencyRate(currency)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(response)
}

func handlePairRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	pair := r.URL.Query().Get("pair")
	response, err := apiClient.GetExchangeRatePairs(pair)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	json.NewEncoder(w).Encode(response)
}

// JavaScript-bound functions
func getAllRates() string {
	response, err := apiClient.GetAllExchangeRates()
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	return formatAllRates(response)
}

func getCurrencyRate(currency string) string {
	if currency == "" {
		return "Error: Please enter a currency code"
	}
	response, err := apiClient.GetCurrencyRate(currency)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	return formatCurrencyRate(response, currency)
}

func getPairRate(pairs string) string {
	if pairs == "" {
		return "Error: Please enter currency pair(s)"
	}
	response, err := apiClient.GetExchangeRatePairs(pairs)
	if err != nil {
		return fmt.Sprintf("Error: %s", err.Error())
	}
	return formatPairRate(response)
}

// Format functions
func formatAllRates(response *api.AllRatesResponse) string {
	if response == nil {
		return "No data available"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== Exchange Rates ===\n"))
	sb.WriteString(fmt.Sprintf("Last Updated: %d\n\n", response.UpdatedAt))

	// Display all currencies and their rates
	for _, entry := range response.Rates {
		sb.WriteString(fmt.Sprintf("--- %s ---\n", entry.Base))

		// Sort currency codes for consistent display
		currencies := make([]string, 0, len(entry.Rates))
		for curr := range entry.Rates {
			currencies = append(currencies, curr)
		}
		sort.Strings(currencies)

		// Display each rate
		for _, curr := range currencies {
			rate := entry.Rates[curr]
			sb.WriteString(fmt.Sprintf("  %s: %.8f\n", curr, rate))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func formatCurrencyRate(response *api.CurrencyResponse, currency string) string {
	if response == nil {
		return "No data available"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("=== %s Exchange Rates ===\n\n", strings.ToUpper(currency)))

	keys := make([]string, 0, len(response.Data))
	for k := range response.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := response.Data[key]
		sb.WriteString(fmt.Sprintf("%s: %v\n", key, formatValue(value)))
	}

	return sb.String()
}

func formatPairRate(response *api.PairResponse) string {
	if response == nil {
		return "No data available"
	}

	var sb strings.Builder
	sb.WriteString("=== Currency Pair Rates ===\n\n")

	keys := make([]string, 0, len(response.Data))
	for k := range response.Data {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := response.Data[key]
		sb.WriteString(fmt.Sprintf("%s: %v\n", key, formatValue(value)))
	}

	return sb.String()
}

func formatValue(value interface{}) string {
	switch v := value.(type) {
	case map[string]interface{}:
		jsonBytes, err := json.MarshalIndent(v, "  ", "  ")
		if err != nil {
			return fmt.Sprintf("%v", v)
		}
		return "\n  " + string(jsonBytes)
	case float64:
		return fmt.Sprintf("%.4f", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func getHTMLContent() string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Exchange Rate App</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }
        
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 50%, #0f3460 100%);
            color: #e8e8e8;
            min-height: 100vh;
        }
        
        .container {
            max-width: 900px;
            margin: 0 auto;
            padding: 20px;
        }
        
        h1 {
            text-align: center;
            margin-bottom: 30px;
            font-size: 2.2rem;
            background: linear-gradient(90deg, #00d4ff, #7c3aed);
            -webkit-background-clip: text;
            -webkit-text-fill-color: transparent;
            background-clip: text;
        }
        
        .tabs {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
        }
        
        .tab-btn {
            flex: 1;
            padding: 15px 20px;
            border: none;
            background: rgba(255, 255, 255, 0.1);
            color: #e8e8e8;
            font-size: 1rem;
            cursor: pointer;
            border-radius: 10px;
            transition: all 0.3s ease;
        }
        
        .tab-btn:hover {
            background: rgba(255, 255, 255, 0.2);
        }
        
        .tab-btn.active {
            background: linear-gradient(135deg, #00d4ff, #7c3aed);
            font-weight: bold;
        }
        
        .tab-content {
            display: none;
            background: rgba(255, 255, 255, 0.05);
            border-radius: 15px;
            padding: 25px;
            border: 1px solid rgba(255, 255, 255, 0.1);
        }
        
        .tab-content.active {
            display: block;
            animation: fadeIn 0.3s ease;
        }
        
        @keyframes fadeIn {
            from { opacity: 0; transform: translateY(10px); }
            to { opacity: 1; transform: translateY(0); }
        }
        
        .form-group {
            margin-bottom: 20px;
        }
        
        label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: #a8a8a8;
        }
        
        input {
            width: 100%;
            padding: 12px 15px;
            border: 2px solid rgba(255, 255, 255, 0.1);
            border-radius: 8px;
            background: rgba(0, 0, 0, 0.2);
            color: #e8e8e8;
            font-size: 1rem;
            transition: border-color 0.3s ease;
        }
        
        input:focus {
            outline: none;
            border-color: #00d4ff;
        }
        
        input::placeholder {
            color: #666;
        }
        
        .btn {
            padding: 12px 30px;
            border: none;
            border-radius: 8px;
            font-size: 1rem;
            cursor: pointer;
            transition: all 0.3s ease;
        }
        
        .btn-primary {
            background: linear-gradient(135deg, #00d4ff, #7c3aed);
            color: white;
        }
        
        .btn-primary:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 20px rgba(0, 212, 255, 0.4);
        }
        
        .btn:disabled {
            opacity: 0.5;
            cursor: not-allowed;
            transform: none;
        }
        
        .result-box {
            margin-top: 20px;
            padding: 20px;
            background: rgba(0, 0, 0, 0.3);
            border-radius: 10px;
            min-height: 200px;
            max-height: 350px;
            overflow-y: auto;
            font-family: 'Consolas', 'Courier New', monospace;
            font-size: 0.9rem;
            line-height: 1.6;
            white-space: pre-wrap;
        }
        
        .result-box::-webkit-scrollbar {
            width: 8px;
        }
        
        .result-box::-webkit-scrollbar-track {
            background: rgba(0, 0, 0, 0.2);
            border-radius: 4px;
        }
        
        .result-box::-webkit-scrollbar-thumb {
            background: rgba(255, 255, 255, 0.2);
            border-radius: 4px;
        }
        
        .loading {
            display: none;
            color: #00d4ff;
            margin-left: 10px;
        }
        
        .loading.show {
            display: inline-block;
        }
        
        .tip {
            font-size: 0.85rem;
            color: #888;
            margin-top: 8px;
        }
        
        .error {
            color: #ff6b6b;
        }
        
        @keyframes spin {
            to { transform: rotate(360deg); }
        }
        
        .spinner {
            display: inline-block;
            width: 16px;
            height: 16px;
            border: 2px solid rgba(0, 212, 255, 0.3);
            border-top-color: #00d4ff;
            border-radius: 50%;
            animation: spin 0.8s linear infinite;
            margin-right: 8px;
            vertical-align: middle;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Exchange Rate App</h1>
        
        <div class="tabs">
            <button class="tab-btn active" onclick="switchTab(this, 'all-rates')">All Rates</button>
            <button class="tab-btn" onclick="switchTab(this, 'currency')">Currency Lookup</button>
            <button class="tab-btn" onclick="switchTab(this, 'pair')">Pair Exchange</button>
        </div>
        
        <div id="all-rates" class="tab-content active">
            <p style="margin-bottom: 20px; color: #a8a8a8;">View all available exchange rates from the API</p>
            <button class="btn btn-primary" onclick="fetchAllRates()" id="allRatesBtn">
                Refresh Rates
            </button>
            <span class="loading" id="allRatesLoading"><span class="spinner"></span>Loading...</span>
            <div class="result-box" id="allRatesResult">Click "Refresh Rates" to load exchange rates...</div>
        </div>
        
        <div id="currency" class="tab-content">
            <div class="form-group">
                <label>Currency Code</label>
                <input type="text" id="currencyInput" placeholder="Enter currency code (e.g., USD, EUR, IDR)" maxlength="3">
            </div>
            <button class="btn btn-primary" onclick="fetchCurrencyRate()" id="currencyBtn">
                Get Rate
            </button>
            <span class="loading" id="currencyLoading"><span class="spinner"></span>Loading...</span>
            <div class="result-box" id="currencyResult">Enter a currency code and click "Get Rate"...</div>
        </div>
        
        <div id="pair" class="tab-content">
            <div class="form-group">
                <label>Currency Pair(s)</label>
                <input type="text" id="pairInput" placeholder="e.g., USDIDR or USDIDR,EURIDR,MYRIDR">
                <p class="tip">Tip: Use comma for multiple pairs (e.g., USDIDR,EURIDR)</p>
            </div>
            <button class="btn btn-primary" onclick="fetchPairRate()" id="pairBtn">
                Get Rate
            </button>
            <span class="loading" id="pairLoading"><span class="spinner"></span>Loading...</span>
            <div class="result-box" id="pairResult">Enter currency pair(s) and click "Get Rate"...</div>
        </div>
    </div>
    
    <script>
        function switchTab(btn, tabId) {
            document.querySelectorAll('.tab-btn').forEach(function(b) { b.classList.remove('active'); });
            document.querySelectorAll('.tab-content').forEach(function(c) { c.classList.remove('active'); });
            
            btn.classList.add('active');
            document.getElementById(tabId).classList.add('active');
        }
        
        async function fetchAllRates() {
            var btn = document.getElementById('allRatesBtn');
            var loading = document.getElementById('allRatesLoading');
            var result = document.getElementById('allRatesResult');
            
            btn.disabled = true;
            loading.classList.add('show');
            result.textContent = 'Fetching data...';
            
            try {
                var response = await window.getAllRates();
                result.textContent = response;
            } catch (err) {
                result.innerHTML = '<span class="error">Error: ' + (err.message || err) + '</span>';
            } finally {
                btn.disabled = false;
                loading.classList.remove('show');
            }
        }
        
        async function fetchCurrencyRate() {
            var btn = document.getElementById('currencyBtn');
            var loading = document.getElementById('currencyLoading');
            var result = document.getElementById('currencyResult');
            var input = document.getElementById('currencyInput').value.trim().toUpperCase();
            
            if (!input) {
                result.innerHTML = '<span class="error">Please enter a currency code</span>';
                return;
            }
            
            btn.disabled = true;
            loading.classList.add('show');
            result.textContent = 'Fetching data...';
            
            try {
                var response = await window.getCurrencyRate(input);
                result.textContent = response;
            } catch (err) {
                result.innerHTML = '<span class="error">Error: ' + (err.message || err) + '</span>';
            } finally {
                btn.disabled = false;
                loading.classList.remove('show');
            }
        }
        
        async function fetchPairRate() {
            var btn = document.getElementById('pairBtn');
            var loading = document.getElementById('pairLoading');
            var result = document.getElementById('pairResult');
            var input = document.getElementById('pairInput').value.trim().toUpperCase();
            
            if (!input) {
                result.innerHTML = '<span class="error">Please enter currency pair(s)</span>';
                return;
            }
            
            btn.disabled = true;
            loading.classList.add('show');
            result.textContent = 'Fetching data...';
            
            try {
                var response = await window.getPairRate(input);
                result.textContent = response;
            } catch (err) {
                result.innerHTML = '<span class="error">Error: ' + (err.message || err) + '</span>';
            } finally {
                btn.disabled = false;
                loading.classList.remove('show');
            }
        }
        
        // Allow Enter key to submit
        document.getElementById('currencyInput').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') fetchCurrencyRate();
        });
        
        document.getElementById('pairInput').addEventListener('keypress', function(e) {
            if (e.key === 'Enter') fetchPairRate();
        });
    </script>
</body>
</html>`
}
