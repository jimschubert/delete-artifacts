FROM gcr.io/distroless/static-debian12
ARG APP_NAME
COPY /${APP_NAME} /app
ENTRYPOINT ["/app"]
