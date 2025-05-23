FROM golang:1.24-alpine AS build

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.7 AS runtime

WORKDIR /app

COPY --from=build /go/src/app/church-members-api .

CMD ["/app/church-members-api"]
