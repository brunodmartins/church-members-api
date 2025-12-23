dev :
	echo "Rump up docker containers (LocalStack)"
	mkdir -p ./volume
	chmod 777 ./volume
	docker compose -f docker-compose.yml -p church-members-api up -d
	echo "Load test dataset"
	bash runbook/create_local_data.sh