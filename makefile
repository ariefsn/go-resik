upgrade.dlv:
	go install github.com/go-delve/delve/cmd/dlv@latest

install.air:
	go install github.com/cosmtrek/air@latest

install.deps:
	go mod tidy

dev.hot:
	air -c .air.toml