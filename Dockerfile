# ---- Builder ----
FROM golang:1.24.4-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
ENV GOPROXY=https://proxy.golang.org,direct
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
    go build -trimpath -ldflags "-s -w" -o /app /src/cmd/app

# ---- Runtime ----
FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /
COPY --from=builder /app /app
EXPOSE 3000
ENTRYPOINT ["/app"]
