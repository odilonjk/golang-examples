docker_build:
	@docker build . -t rabbitmq-delayed -q
	@echo 'RabbitMQ image with Delayed Message plugin built successfully'

rabbit_start:
	@docker run -d --rm --name rabbitmq -p 15672:15672 -p 5672:5672 rabbitmq-delayed:latest
	@echo 'RabbitMQ started'
	@echo 'Management available on http://localhost:15672'

rabbit_stop:
	@docker stop rabbitmq
	@echo 'RabbitMQ shutted down'

estoque_run:
	@cd estoque && go run . && cd ..

ordem_producao_start:
	@cd ordem-de-producao && go run .
