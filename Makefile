up: ##build and run project in docker conteiner
	@docker-compose up --build -d




test: ##run test
	@go test -v ./...
    


kvput: ##put value into  consul kv store
	@./config/register-variables.sh
	
down: ##stop and remobe all conteiner
	@docker-compose down --remove-orphans
