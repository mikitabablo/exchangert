FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o exchangert ./cmd/http/main.go

FROM gcr.io/distroless/base-debian12

WORKDIR /app

COPY --from=builder /app/exchangert .

EXPOSE 8080

ENTRYPOINT ["./exchangert"]
