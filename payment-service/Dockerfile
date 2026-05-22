FROM golang:1.25-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/payment-service ./payment-service

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/payment-service /usr/local/bin/payment-service

EXPOSE 8080
USER nonroot:nonroot
ENTRYPOINT ["/usr/local/bin/payment-service"]
