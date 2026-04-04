# Stage 1: Build frontend
FROM node:20-alpine AS frontend
WORKDIR /app/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build backend
FROM golang:1.26-alpine AS backend
WORKDIR /app
COPY backend/go.* ./
RUN go mod download
COPY backend/ ./
COPY --from=frontend /app/frontend/dist ./static
RUN CGO_ENABLED=0 go build -o server ./cmd/server

# Stage 3: Runtime
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=America/Sao_Paulo
WORKDIR /app
COPY --from=backend /app/server .
COPY --from=backend /app/static ./static
COPY --from=backend /app/migrations ./migrations
COPY --from=backend /app/seed ./seed
EXPOSE 8081
CMD ["./server"]
