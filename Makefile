run_db:
	go run ./cmd/migrator --migrations-path=./migrations 

run: 
	go run cmd/sso/main.go --config=./config/local.yaml
	