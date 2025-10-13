migrate :
	migrate -source file://database/migrations \ 
	-database postgresql://root:root@localhost:5432/go_sweeper_dev_db?sslmode=disable up 

rollback :
	migrate -source file://database/migrations \ 
	-database postgresql://root:root@localhost:5432/go_sweeper_dev_db?sslmode=disable down 

drop :
	migrate -source file://database/migrations \ 
	-database postgresql://root:root@localhost:5432/go_sweeper_dev_db?sslmode=disable drop 
	
migration : 
	@read -p "Enter name of migration: " name; \
		migrate create -ext sql -dir database/migrations $$name

run:
	go run cmd/graphqlserver/main.go