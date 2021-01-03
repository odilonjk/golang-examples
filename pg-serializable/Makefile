build:
	@echo "Gerando binario da aplicacao"
	@go build -o booking_app

build_docker:
	@docker build . -t booking-app

start:
	@docker-compose up --detach

stop:
	@docker-compose down

concurrent_calls:
	@go run ./cmd/concurrent.go
	@echo Done!