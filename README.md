# Country and Exchange API

## About
This API enables a user to retrieve information about countries, as well as the exchange rates of the neighbouring countries.

## Endpoints
The API contains the following endpoints.
```text
/countryinfo/v1/status/ 
/countryinfo/v1/info/{country}
/countryinfo/v1/exchange/{country}
```

### Endpoint Usage


## Dependencies
The API is dependent on two upstream API's. 
```text
Country API: http://129.241.150.113:8080/v3.1/
Currency API: http://129.241.150.113:9090/currency/
```

## AI-Usage