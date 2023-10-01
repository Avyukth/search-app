
FROM golang:1.21 as builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o search-api ./cmd/server/main.go

FROM gcr.io/distroless/static-debian11
COPY --from=builder /app/search-api /search-api
CMD ["/search-api"]
