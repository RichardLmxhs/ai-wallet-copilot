.PHONY: build run test clean docker-build docker-run

build:
	go build -o bin/ai-wallet-copilot ./cmd/server

run: build
	./bin/ai-wallet-copilot

test:
	go test ./...

clean:
	rm -rf bin/

docker-build:
	docker build -t ai-wallet-copilot -f deployments/docker/Dockerfile .

docker-run: docker-build
	docker run -p 8080:8080 ai-wallet-copilot
