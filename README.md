# Gome

Gome is a personal dashboard I use in the morning to get useful information such as next tram departures or the weather.

## Usage

First, a token for the weather api should be created from https://developer.forecast.io/

Then, there is no need to clone the repository, a Docker image is shipped straight to the Docker Hub :

```
docker run -p 8080:80 -e "FORECAST_API_KEY=my_api_token" bamarni/gome:latest
```

## Todo

* integration with the BVG API (http://www.vbb.de/de/article/webservices/schnittstellen-fuer-webentwickler/5070.html)
