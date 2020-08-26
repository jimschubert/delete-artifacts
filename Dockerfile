FROM golang:1.14-alpine3.11 as builder
ARG APP_NAME

ENV GOOS=linux \
    GOARCH=386 \
    CGO_ENABLED=0

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download && go build -o /go/bin/app ./cmd

FROM gcr.io/distroless/base-debian10
COPY --from=builder /go/bin/app /
ENTRYPOINT ["/app"]
