## Requirements

- Write a service in go-lang that exposes a web API that, when called, will tell you whether or not the International Space Station is currently visible to those on the earth directly beneath it.

- The space station should be considered visible when the following conditions are met:
  - a) It is after sunset, and
  - b) cloud cover <= 30%

> Useful public APIs to use:
>
> To determine the position of the space station at any instant you can use the API documented here: http://open-notify.org/Open-Notify-API/ISS-Location-Now/
>
> To determine local time, sunrise/sunset, and cloud cover, you can use the API documented here: https://www.weatherbit.io/api/weather-current
>
> Note you can sign up for a free key to use this API – it will be limited to 500 calls per day, at a max rate of 1 call/second, which should be enough for this exercise.
>
> Instructions:
>
> Please write the service as you would your production-quality code, including documentation, unit tests, and other measures you would take to ensure the correctness of the code. Also have a think about what V2 of the application might look like.

### TODO

- [x] Add ISS location service
- [x] Add Weather service
- [x] Call services in sequence
- [x] Call services in parallel
- [x] Replace default router with chi

### MVP

- [x] Configurable accuracy of overhead via ENV Variable
  - [ ] via (optional) query param
- [x] function composition of handlers
- [x] allow setting accuracy from ENV variable
- [ ] Handle non-200 responses from external rest apis
- [ ] gzip compress response
- [x] generate swagger docs
- [ ] structured logging
  - [x] set log level via ENV Variable

### Additional

- [ ] serve over https
- [ ] Add dev `Dockerfile`
- [ ] Add build `Dockerfile`
- [ ] static code analysis (`go vet`/`staticcheck`)
- [ ] serve swagger-ui
- [ ] Circle-ci build
- [ ] Code coverage report
- [ ] middleware to log requests/responses
