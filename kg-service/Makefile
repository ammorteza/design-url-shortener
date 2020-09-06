db_migrate:
	sudo docker exec kgs_server1 go build -o cmd/main cmd/main.go
	sudo docker exec kgs_server1 ./cmd/main db:migrate

db_reset:
	sudo docker exec kgs_server1 go build -o cmd/main cmd/main.go
	sudo docker exec kgs_server1 ./cmd/main db:reset

db_seed:
	sudo docker exec kgs_server1 go build -o cmd/main cmd/main.go
	sudo docker exec kgs_server1 ./cmd/main db:seed