up: ##build and run project in docker conteiner
	@docker-compose up --build -d
test: ##run test
	  ##This task for future		
down: ##stop and remobe all conteiner
	@docker-compose down --remove-orphans