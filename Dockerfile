FROM golang:1.13 AS build

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && go install -v ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.7 AS runtime

WORKDIR /app

ENV VPR_CHURCH_MEMBERS_DATABASE_URL $VPR_CHURCH_MEMBERS_DATABASE_URL
ENV VPR_CHURCH_MEMBERS_ACCESS_TOKEN $VPR_CHURCH_MEMBERS_ACCESS_TOKEN
ENV SCOPE $SCOPE

COPY --from=build /go/src/app/config/* ./config/
COPY --from=build /go/src/app/church-members-api .

CMD ["/app/church-members-api"]