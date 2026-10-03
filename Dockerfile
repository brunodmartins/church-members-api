# syntax=docker/dockerfile:1.7

FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -trimpath -ldflags="-s -w" -o /out/church-members-api .

FROM alpine:3.7 AS runtime

WORKDIR /app

COPY --from=build /out/church-members-api .

CMD ["/app/church-members-api"]
