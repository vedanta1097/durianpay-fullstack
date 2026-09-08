FROM golang:1.21-bookworm AS build

RUN apt-get update \
    && apt-get install -y --no-install-recommends build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=1 go build -o /out/durianpay-api main.go

FROM debian:bookworm-slim

RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /out/durianpay-api ./durianpay-api
COPY openapi.yaml ./openapi.yaml

ENV HTTP_ADDR=:8080
ENV DB_PATH=/data/dashboard.db
ENV OPENAPIYAML_LOCATION=/app/openapi.yaml

EXPOSE 8080

CMD ["./durianpay-api"]
