Aplikasi Go berbasis GUI

Real-Time Exchange Rate API 

library : Fyne → GUI native, simpel

env : api_key = ""

Base URL = https://use.api.co.id

Endpoint :
GET /api/exchange-rates
    Parameters : 
    - x-api-co-id (string)
GET /currency/:currency
    Parameters : 
    - x-api-co-id (string)
    - currency *  (string) (path) Base currency code (3 characters, case insensitive, e.g., USD, EUR, IDR)
GET /currency/exchange-rate
    Parameters : 
    - x-api-co-id (string)
    - pair * (string) (query) Currency pair symbol(s). Single pair (e.g., USDIDR) or comma-separated multiple pairs (e.g., USDIDR,EURIDR,MYRIDR)



