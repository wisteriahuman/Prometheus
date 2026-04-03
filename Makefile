.PHONY: build web clean dev install release

web:
	cd web && pnpm install && pnpm run build

build: web
	go build -o prm ./cmd/prometheus

dev:
	go run ./cmd/prometheus dev $(VAULT) -p $(or $(PORT),3000)

install: web
	go install ./cmd/prometheus

clean:
	rm -rf prm internal/server/web
	cd web && rm -rf build .svelte-kit

release: web
	GOOS=darwin GOARCH=arm64 go build -o prm-darwin-arm64 ./cmd/prometheus
	GOOS=linux GOARCH=amd64 go build -o prm-linux-amd64 ./cmd/prometheus
