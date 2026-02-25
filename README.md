# Country and Exchange API

## About
This API enables a user to retrieve information about countries, as well as the exchange rates of the neighbouring countries.

## Usage
```
go build assignment_one/cmd/app
./app
```

## Endpoints
The API contains the following endpoints. Country has to be on the ISO3166-1 alpha-2 format.
```text
/countryinfo/v1/status/ 
/countryinfo/v1/info/{country}
/countryinfo/v1/exchange/{country}
```

## Dependencies
The API is dependent on two upstream API's. 
```text
Country API: http://129.241.150.113:8080/v3.1/
Currency API: http://129.241.150.113:9090/currency/
```