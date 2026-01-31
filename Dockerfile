FROM golang:1.25-alpine as builder
ARG APP_NAME

ENV GOOS=linux \
    GOARCH=386 \
    CGO_ENABLED=0

WORKDIR /go/src/app
ADD . /go/src/app

RUN go mod download && go build -o /go/bin/app ./cmd

FROM gcr.io/distroless/static-debian12
COPY --from=builder /go/bin/app /
ENTRYPOINT ["/app"]
