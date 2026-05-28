
FROM oven/bun:1 AS ui-builder
WORKDIR /ui
COPY ui/package.json ui/bun.lock ./
RUN bun install --frozen-lockfile
COPY ui/ ./
RUN bun run build

FROM golang:1.26-alpine AS go-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui-builder /ui/dist ./ui/dist
RUN CGO_ENABLED=0 GOOS=linux go build -o vivery ./cmd/vivery
RUN go install github.com/jackc/tern/v2@latest

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=go-builder /app/vivery ./vivery
COPY --from=go-builder /go/bin/tern ./tern
COPY migrations/ ./migrations/
COPY docker-entrypoint.sh ./docker-entrypoint.sh
RUN chmod +x ./docker-entrypoint.sh
EXPOSE 3000
ENTRYPOINT ["./docker-entrypoint.sh"]
