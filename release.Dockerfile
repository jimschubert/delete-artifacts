FROM gcr.io/distroless/base-debian10
ARG APP_NAME
COPY /${APP_NAME} /app
ENTRYPOINT ["/app"]
