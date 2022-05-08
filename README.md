# the-starry-night

### Build

1. To build executable run `make` from the root directory
1. To run the compiled binary run `./main`

|                                          | command      |
| ---------------------------------------- | ------------ |
| Build then run app                       | `make start` |
| Watch files then auto-compile then start | `make dev`   |

### Configuration

The default values are loaded from `settings.yaml` and can be overridden by environment variables

| Description                     | ENV VARIABLE  | DEFAULT VALUE |
| ------------------------------- | ------------- | ------------- |
| Port that application listen on | `SERVER_PORT` | `8080`        |

# Prerequisites

- `go` version `1.17`+
- `go-swagger` [installation](https://goswagger.io/install.html)
- Recommended
  - For live-reload during development `air` to install run: `go install github.com/cosmtrek/air@latest`
