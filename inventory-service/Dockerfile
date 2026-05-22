FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/inventory-service ./inventory-service

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/inventory-service /usr/local/bin/inventory-service

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/inventory-service"]
