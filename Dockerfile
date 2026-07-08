FROM gcr.io/distroless/static-debian12:nonroot

COPY --chown=nonroot:nonroot dist/cooperative-api /api

EXPOSE 8089
ENTRYPOINT ["/api"]
