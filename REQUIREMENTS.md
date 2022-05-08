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
- [ ] Add Weather service
- [ ] Replace default router with chi

### MVP

- [ ] Configurable accuracy of overhead via ENV Variable
  - [ ] via query param
- [ ] function composition of handlers
- [ ] allow setting accuracy from ENV variable
- [ ] Handle non-200 responses from external rest apis
- [ ] gzip compress response
- [ ] generate swagger docs
- [ ] structured logging
  - [ ] set log level ENV Variable

### Additional

- [ ] serve over https
- [ ] dev `Dockerfile`
- [ ] Add `Dockerfile`
- [ ] static code analysis (`go vet`/`staticcheck`)
- [ ] serve swagger-ui
- [ ] Circle-ci build
- [ ] Code coverage report
- [ ] middleware to log requests/responses
