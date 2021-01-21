FROM golang:1.15-alpine AS build
WORKDIR /src
COPY . .
RUN go build .

FROM alpine
COPY --from=build /src/play.spruce.cf /usr/bin/play

WORKDIR /app
COPY ./assets /app/assets

EXPOSE 8080
CMD ["play"]
