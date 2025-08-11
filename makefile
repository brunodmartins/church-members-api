dev :
	echo "Rump up docker containers (LocalStack)"
	mkdir -p ./volume
	chmod 777 ./volume
	docker compose -f docker-compose.yml -p church-members-api up -d
	echo "Building infrastructure with OpenTofu"
	tofu -chdir=terraform/environment/local apply -auto-approve
	echo "Load test dataset"
	bash runbook/create_local_data.sh