upgrade.dlv:
	go install github.com/go-delve/delve/cmd/dlv@latest

install.air:
	go install github.com/cosmtrek/air@latest

install.deps:
	go mod tidy

dev.hot:
	air -c .air.toml

mock.repository:
	cd ./domain && mockery \
	--name ${app}Repository \
	--filename $(shell echo ${app} | tr '[:upper:]' '[:lower:]')_repository.go \
	--outpkg mocks \
	--structname "${app}Repository"

mock.service:
	cd ./domain && mockery \
	--name ${app}Service \
	--filename $(shell echo ${app} | tr '[:upper:]' '[:lower:]')_service.go \
	--outpkg mocks \
	--structname "${app}Service"

run:
	go run main.go

build:
	go build -o main .

build.run:
	go build -o main . && ./main

test.coverage.html:
	go test ./app/... -coverprofile=cover.out && go tool cover -html=cover.out -o coverage.html

test.coverage.func:
	go test ./app/... -coverprofile=cover.out && go tool cover -func=cover.out -o coverage.html

test.coverage:
	go test ./app/... -race -v -coverprofile=cover.out && ./scripts/coverage.sh cover.out ${threshold}