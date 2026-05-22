FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/order-service ./order-service

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/order-service /usr/local/bin/order-service

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/order-service"]
