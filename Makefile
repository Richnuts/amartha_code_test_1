include .env

create.migration:
	migrate create -ext sql -dir ./migrations -seq ${NAME}

migrate.migration:
	migrate -database "postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:${POSTGRES_PORT}/${POSTGRES_DATABASE}?sslmode=disable" -path migrations up ${N}

### Install Dependency
dependency:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go get -u github.com/labstack/echo/v4
	go get -u github.com/doug-martin/goqu/v9
	go get github.com/robfig/cron/v3@v3.0.0
	go mod tidy