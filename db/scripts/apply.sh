export GOOSE_DRIVER=postgres
export GOOSE_DBSTRING="host=127.0.0.1 user=postgres password=postgres dbname=postgres port=5432 sslmode=disable"

goose -dir "./../migrations" up