# === BCE Belgium ===

bce-up:
	docker compose -f bce_belgium_backend/docker-compose.yml up -d

bce-down:
	docker compose -f bce_belgium_backend/docker-compose.yml down

bce-logs:
	docker compose -f bce_belgium_backend/docker-compose.yml logs -f

bce-build:
	cd bce_belgium_backend && go build -o bce-api .

bce-api:
	cd bce_belgium_backend && go run . api

bce-import:
	cd bce_belgium_backend && go run . all

# === SIRENE France ===

sirene-up:
	docker compose -f sirene_france_backend/docker-compose.yml up -d

sirene-down:
	docker compose -f sirene_france_backend/docker-compose.yml down

sirene-logs:
	docker compose -f sirene_france_backend/docker-compose.yml logs -f

sirene-build:
	cd sirene_france_backend && go build -o sirene-api .

sirene-api:
	cd sirene_france_backend && go run . api

sirene-import:
	cd sirene_france_backend && go run . all

# === Both ===

up-all:
	docker compose -f bce_belgium_backend/docker-compose.yml up -d
	docker compose -f sirene_france_backend/docker-compose.yml up -d

down-all:
	docker compose -f bce_belgium_backend/docker-compose.yml down
	docker compose -f sirene_france_backend/docker-compose.yml down

ps:
	@echo "=== BCE Belgium ==="
	@docker compose -f bce_belgium_backend/docker-compose.yml ps
	@echo ""
	@echo "=== SIRENE France ==="
	@docker compose -f sirene_france_backend/docker-compose.yml ps
