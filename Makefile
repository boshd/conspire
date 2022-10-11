build:
	docker-compose build

run:
	docker-compose up -d --build

restart:
	make stop && make build && make run

stop:
	docker-compose down -v