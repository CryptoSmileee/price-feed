
<p align="center">
  <a href="https://orionprotocol.io">
    <img src="https://orionprotocol.io/images/logo.png" />
  </a>
</p>

<h3 align="center">Orion Price Feed API Documentation</h3>

## API Functions

- [API Functions](#api-functions)
- [Order book](#order-book)
  - [HTTP Request](#http-request)
  - [Query Parameters](#query-parameters)
- [Candles](#candles)
  - [HTTP Request](#http-request-1)
  - [Query Parameters](#query-parameters-1)
- [Reload](#reload)
  - [HTTP Request](#http-request-2)
  - [Query Parameters](#query-parameters-2)

## Order book

This endpoint retrieves the order book

### HTTP Request

`GET http://example.com/api/v1/orderBook`

### Query Parameters

Parameter | Description
-|-
symbol | Trading pair for which order book is requested
depth | Oder book depth between 1 and 1000

## Candles

This endpoint retrieves the candles for a trading pair

### HTTP Request

`GET http://example.com/api/v1/candles`

### Query Parameters

Parameter | Description
-|-
symbol | Trading pair for which the candles are requested
depth | Oder book depth between 1 and 1000

## Reload

Reload the Api

### HTTP Request

`GET http://example.com/api/v1/reload`

### Query Parameters

Parameter | Description
-|-
token | Api token specified in config.json