up:
	docker-compose up -d

down:
	docker-compose down

exec:
	docker exec -it cms-api_postgres psql -U admin -d cms_api


 go test ./tests/..