docker-up:
	docker compose up --build -d

docker-down:
	docker compose down

migrate-up:
	docker compose run migrate up 

migrate-down:
	docker compose run migrate down

test-service:
	cd internal/service && go test -v
