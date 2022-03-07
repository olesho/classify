FROM golang:1.15-alpine AS dev

ENV APP_NAME="main" APP_PATH="/var/app"
COPY . ${APP_PATH}
WORKDIR ${APP_PATH}
ENV CGO_ENABLED=0 GOOS=linux
RUN apk update && apk add git
RUN go get "github.com/olesho/classify/sequence"
RUN go build -o main
RUN chmod +x main

FROM scratch AS prod
ENV APP_BUILD_PATH="/var/app" \
    APP_BUILD_NAME="main"
WORKDIR ${APP_BUILD_PATH}
COPY --from=dev ${APP_BUILD_PATH}/${APP_BUILD_NAME} ${APP_BUILD_PATH}/
EXPOSE 9876
ENTRYPOINT ["/var/app/main"]
CMD ""