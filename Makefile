.PHONY: s, c

s:
	@go run ./cmd/tictactoe/

c:
	@cd web
	@npm run dev

# csql --url $DATABASE_URL -f ./ccdb/migration.sql 