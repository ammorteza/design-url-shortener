db_migrate:
	sudo docker exec app_server1 go build -o cmd/main cmd/main.go
	sudo docker exec app_server1 ./cmd/main db:migrate

db_reset:
	sudo docker exec app_server1 go build -o cmd/main cmd/main.go
	sudo docker exec app_server1 ./cmd/main db:reset

test:
	sudo docker exec app_server1 go test . -count=1000 -timeout 40m
