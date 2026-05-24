.PHONY: rebuild logs shell down

rebuild:
	docker compose down -v
	docker compose up -d --build

rebuild-cache:
	docker compose down
	docker compose up -d --build

down:
	docker compose down -v

logs:
	docker compose logs -f
