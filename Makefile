.PHONY: dev dev-backend dev-frontend build up down seed clean

dev-backend:
	cd backend && go run ./cmd/server

dev-frontend:
	cd frontend && npm run dev

dev:
	@echo "Execute em terminais separados:"
	@echo "  make dev-backend"
	@echo "  make dev-frontend"

build:
	docker compose build

up:
	docker compose up -d

down:
	docker compose down

seed:
	cd backend && go run ./cmd/server --seed

clean:
	docker compose down -v
	rm -f data/mercado.db
