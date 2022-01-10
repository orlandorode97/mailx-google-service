build:
	docker-compose build --no-cache

run:
	docker-compose up

build-run:
	docker-compose up --build --remove-orphans