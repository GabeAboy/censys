up:
	docker compose up -d

down:
	docker compose down

build:
	docker compose build

restart:
	docker compose restart

rebuild:
	docker compose down
	docker compose build --no-cache
	docker compose up -d

clean: ## Remove all containers, volumes, and images
	docker compose down -v --rmi all

db-shell:
	docker compose exec postgres psql -U censys -d censys

api-shell:
	docker compose exec api sh

logs:
	docker compose logs -f

logs-api:
	docker compose logs -f api

logs-db:
	docker compose logs -f postgres

up-postgres:
	docker compose up -d postgres

swagger-gen:
	swag init -g main.go --output docs