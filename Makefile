stop:
	sudo systemctl stop postgresql && sudo systemctl stop redis.service
start:
	docker-compose up -d
stats:
	docker stats
