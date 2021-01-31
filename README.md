# kahvibot

Kahvikanavalle tehty turhake

## Käyttö

- PostgreSQL pystyyn ja `kahvidb.sql` sisään.
- Kantaan käyttäjä botille, oikeudet tauluihin.
- Luo `.env` tiedosto missä nämä kentät:
```
TG_TOKEN=
PSQL_URL= (connection URI)
```
- `go run main.go`
