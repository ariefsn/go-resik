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