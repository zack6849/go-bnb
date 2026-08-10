FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy dependency files first so Docker caches this layer
# separately from source changes — deps only re-download when go.mod changes
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# CGO_ENABLED=0 forces a fully static binary — no C library dependencies,
# which is what lets the final image be so minimal
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/seed ./cmd/seed

# ---- Final stage: no Go installed at all, just the binaries ----
FROM alpine:latest

# only needed if your app makes outbound HTTPS calls (real APIs, etc.)
RUN apk add --no-cache ca-certificates

COPY --from=builder /bin/api /bin/api
COPY --from=builder /bin/seed /bin/seed

EXPOSE 8080
CMD ["/bin/api"]