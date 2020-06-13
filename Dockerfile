FROM golang:1.13 AS build

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && go install -v ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.7 AS runtime

WORKDIR /app

ENV VPR_CHURCH_MEMBERS_DATABASE_URL $VPR_CHURCH_MEMBERS_DATABASE_URL
ENV VPR_CHURCH_NAME $VPR_CHURCH_NAME
ENV VPR_AUTH_ENABLE $VPR_AUTH_ENABLE
ENV VPR_AUTH_JWK $VPR_AUTH_JWK
ENV VPR_AUTH_ISS $VPR_AUTH_ISS
ENV VPR_AUTH_TOKEN $VPR_AUTH_TOKEN
ENV VPR_AUTH_AUD $VPR_AUTH_AUD
ENV SCOPE $SCOPE
ENV APP_LANG $APP_LANG

COPY --from=build /go/src/app/config/* ./config/
COPY --from=build /go/src/app/church-members-api .

CMD ["/app/church-members-api"]