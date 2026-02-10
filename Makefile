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

sirene-indexes:
	cd sirene_france_backend && go run . indexes

sirene-reimport:
	docker compose -f sirene_france_backend/docker-compose.yml down -v
	docker compose -f sirene_france_backend/docker-compose.yml up -d
	@echo "Attente du démarrage de PostgreSQL..."
	@sleep 5
	cd sirene_france_backend && go run . all

# === Les deux ===

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

sirene-sql:
	@export $$(grep -v '^#' sirene_france_backend/.env | xargs) && PSQLRC=/dev/null PAGER=cat PGPASSWORD=$$POSTGRES_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$POSTGRES_USER -d $$POSTGRES_DB

sirene-count:
	@export $$(grep -v '^#' sirene_france_backend/.env | xargs) && PSQLRC=/dev/null PAGER=cat PGPASSWORD=$$POSTGRES_PASSWORD psql -h $$DB_HOST -p $$DB_PORT -U $$POSTGRES_USER -d $$POSTGRES_DB -c "SELECT 'unite_legale' as table_name, COUNT(*) FROM unite_legale UNION ALL SELECT 'etablissement', COUNT(*) FROM etablissement;"

help:
	@echo "BCE Belgium:"
	@echo "  make bce-up          Démarrer PostgreSQL + Redis (ports 5433/6379)"
	@echo "  make bce-down        Arrêter les conteneurs"
	@echo "  make bce-logs        Voir les logs"
	@echo "  make bce-build       Compiler le binaire"
	@echo "  make bce-api         Lancer l'API (port 8080)"
	@echo "  make bce-import      Importer les CSV"
	@echo ""
	@echo "SIRENE France:"
	@echo "  make sirene-up       Démarrer PostgreSQL + Redis (ports 5434/6380)"
	@echo "  make sirene-down     Arrêter les conteneurs"
	@echo "  make sirene-logs     Voir les logs"
	@echo "  make sirene-build    Compiler le binaire"
	@echo "  make sirene-api      Lancer l'API (port 8081)"
	@echo "  make sirene-import   Importer les ZIP (+ création indexes)"
	@echo "  make sirene-indexes  Créer les indexes PostgreSQL manuellement"
	@echo "  make sirene-reimport Ré-import complet (supprime volumes + reimporte)"
	@echo "  make sirene-sql      Ouvrir un terminal SQL (sans pager)"
	@echo "  make sirene-count    Compter les lignes dans les tables"
	@echo ""
	@echo "Global:"
	@echo "  make up-all          Démarrer les deux stacks"
	@echo "  make down-all        Arrêter les deux stacks"
	@echo "  make ps              État des conteneurs"
	@echo "  make help            Afficher cette aide"
