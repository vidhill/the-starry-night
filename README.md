# the-starry-night

### Build

1. To build executable run `make` from the root directory
1. To run the compiled binary run `./main`

An api key from [weatherbit.io](https://www.weatherbit.io/api) is necessary to run this application.

This API key should **must** be configured using an environment variable `WEATHER_BIT_API_KEY`

|                                          | command                                        |
| ---------------------------------------- | ---------------------------------------------- |
| Build then run app                       | `WEATHER_BIT_API_KEY=your-api-key; make start` |
| Watch files then auto-compile then start | `WEATHER_BIT_API_KEY=your-api-key; make dev`   |

### Configuration

The default values are loaded from `settings.yaml` and can be overridden by environment variables

| Description                     | ENV VARIABLE                | DEFAULT VALUE                             | Required |
| ------------------------------- | --------------------------- | ----------------------------------------- | -------- |
| Port that application listen on | `SERVER_PORT`               | `8080`                                    |          |
| ISS Rest API URL                | `ISS_API_URL`               | `http://api.open-notify.org/iss-now.json` |          |
| Weatherbit Rest API key         | `WEATHER_BIT_API_KEY`:      | _none_                                    | yes      |
| Weatherbit Rest API base URL    | `WEATHER_BIT_API_BASE_URL`: | `https://api.weatherbit.io/v2.0`          |          |

# Prerequisites

- `go` version `1.17`+
- `go-swagger` [installation](https://goswagger.io/install.html)
- Recommended
  - For live-reload during development `air` to install run: `go install github.com/cosmtrek/air@latest`
