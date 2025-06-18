include .env
export

# Caminho para o binário do Goose
GOOSE=goose

# Caminho para as migrations
MIGRATIONS_DIR=./pkg/database/migrations
TABLE_NAME=migrations

# Driver e DSN de conexão (ajuste conforme necessário)
DRIVER=postgres
DSN=postgres://${DB_USER}:${DB_PASS}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable

# Comando para aplicar todas as migrations
up:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DSN)" up -table ${TABLE_NAME}

# Comando para desfazer a última migration
down:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DSN)" down -table ${TABLE_NAME}

# Status das migrations
status:
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DSN)" status

# Criar nova migration
create:
ifndef NAME
	$(error "Você precisa informar NAME: make create NAME=nome_da_migration")
endif
	$(GOOSE) create $(NAME) -dir $(MIGRATIONS_DIR) -s sql

# Rodar uma migration específica
up-to:
ifndef VERSION
	$(error "Você precisa informar VERSION: make up-to VERSION=20240618120000")
endif
	$(GOOSE) -dir $(MIGRATIONS_DIR) $(DRIVER) "$(DSN)" up-to $(VERSION)
