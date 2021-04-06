FROM golang:1.13 AS build

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && go install -v ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.7 AS runtime

WORKDIR /app

ENV VPR_CHURCH_MEMBERS_DATABASE_URL $VPR_CHURCH_MEMBERS_DATABASE_URL \
    VPR_CHURCH_NAME $VPR_CHURCH_NAME \
    VPR_AUTH_ENABLE $VPR_AUTH_ENABLE \
    VPR_AUTH_JWK $VPR_AUTH_JWK \
    VPR_AUTH_ISS $VPR_AUTH_ISS \
    VPR_AUTH_TOKEN $VPR_AUTH_TOKEN \
    VPR_AUTH_AUD $VPR_AUTH_AUD \
    SCOPE $SCOPE \
    APP_LANG $APP_LANG \
    SERVER $SERVER

COPY --from=build /go/src/app/bundles/* ./bundles/
COPY --from=build /go/src/app/fonts/* ./fonts/
COPY --from=build /go/src/app/church-members-api .

CMD ["/app/church-members-api"]