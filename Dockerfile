# syntax=docker/dockerfile:1

# --- Frontend builder ---
FROM node:22-alpine AS frontend
WORKDIR /app/web
COPY web/package.json web/pnpm-lock.yaml ./
RUN --mount=type=cache,target=/root/.local/share/pnpm/store \
    corepack enable && pnpm install --frozen-lockfile
COPY web/ .
RUN corepack enable && pnpm build

# --- Go builder ---
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
COPY --from=frontend /app/web/dist ./web/dist
ARG VERSION=dev
ARG COMMIT=unknown
RUN --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/server ./cmd/server

# --- Development ---
FROM golang:1.25-alpine AS development
WORKDIR /app
RUN go install github.com/air-verse/air@latest
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
EXPOSE 8080
CMD ["air"]

# --- Production ---
FROM gcr.io/distroless/static-debian12:nonroot AS production
COPY --from=builder /app/server /server
COPY --from=builder /app/web/dist /web/dist
EXPOSE 8080
ENTRYPOINT ["/server"]
