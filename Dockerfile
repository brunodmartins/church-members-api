FROM golang:1.13 AS build

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./... && go install -v ./... && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

FROM alpine:3.7 AS runtime

WORKDIR /app

ENV SERVER $SERVER \
    APPLICATION $APPLICATION \
    CHURCH_NAME $CHURCH_NAME \
    CHURCH_NAME_SHORT $CHURCH_NAME_SHORT \
    APP_LANG $APP_LANG \
    TABLE_MEMBER $TABLE_MEMBER \
    TABLE_MEMBER_HISTORY $TABLE_MEMBER_HISTORY \
    JOBS_DAILY_PHONE $JOBS_DAILY_PHONE \
    REPORTS_TOPIC $REPORTS_TOPIC

COPY --from=build /go/src/app/bundles/* ./bundles/
COPY --from=build /go/src/app/fonts/* ./fonts/
COPY --from=build /go/src/app/church-members-api .

CMD ["/app/church-members-api"]