docker_up:
	@docker-compose up -d

consul_up:
	@./register-services.sh && \
	 ./register-variables.sh


test:
	
up: docker_up consul_up

down:
	@docker-compose down
