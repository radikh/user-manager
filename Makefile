docker_up: ##build and run project in docker conteiner
	@docker-compose up --build -d
test: ##run test
	  ##This task for future		
kvput: ##put value into  consul kv store
	@./infra/consul/register-variables.sh
up: docker_up kvput
down: ##stop and remobe all conteiner
	@docker-compose down --remove-orphans