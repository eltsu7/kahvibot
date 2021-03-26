package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx"
)

func dbKirjaus(userID int, aika int, kuvaus string, nimi string) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Lisää uusi kupillinen
	sql := "insert into juonnit (user_id, aika, kuvaus) values ($1, to_timestamp($2), $3)"

	rows, err := conn.Query(context.Background(), sql, userID, aika, kuvaus)
	if err != nil {
		fmt.Println(rows)
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	rows.Close()

	// Päivitä nimi nimitauluun
	sql = "INSERT INTO nimet VALUES ($1, $2) ON CONFLICT (user_id) DO UPDATE SET username = EXCLUDED.username"

	rows, err = conn.Query(context.Background(), sql, userID, nimi)
	if err != nil {
		fmt.Println(rows)
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	rows.Close()
}

func dbUusinKuppi(userID int) (error, string) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Hae viimeisimmän kirjatun kupin kuvaus
	var kuvaus string
	sql := "SELECT kuvaus FROM juonnit WHERE user_id=$1 ORDER BY id DESC LIMIT 1;"
	err = conn.QueryRow(context.Background(), sql, userID).Scan(&kuvaus)
	if err == pgx.ErrNoRows {
		return errors.New("Ei edellistä kuppia"), ""
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return nil, kuvaus
}

func dbViimeisimmat(userID int, timestamp bool) []string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	sql := "select kuvaus, aika from juonnit where user_id=$1 order by aika desc limit 5"

	rows, _ := conn.Query(context.Background(), sql, userID)

	var viestirivit []string

	for rows.Next() {
		var kuvaus string
		var aika time.Time

		err := rows.Scan(&kuvaus, &aika)

		if err != nil {
			log.Println(err)
			os.Exit(1)
		}

		if timestamp {
			viestirivit = append(viestirivit, fmt.Sprint(aika.Unix())+" "+kuvaus+"\n")
		} else {
			loc, _ := time.LoadLocation("Europe/Helsinki")
			viestirivit = append(viestirivit, aika.In(loc).Format("2.1. 15:04")+" "+kuvaus+"\n")
		}

	}

	return viestirivit
}

func dbPoista(userID int, aika int) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	sql := "delete from juonnit where user_id=$1 and aika=to_timestamp($2)"

	rows, err := conn.Query(context.Background(), sql, userID, aika)
	if err != nil {
		fmt.Println(rows)
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	rows.Close()
}

func dbEiku(uusiKuvaus string, userID int, aika int) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Muuta kuvaus uuteen
	sql := "update juonnit set kuvaus=$1 where user_id=$2 and aika=to_timestamp($3)"

	rows, err := conn.Query(context.Background(), sql, uusiKuvaus, userID, aika)
	if err != nil {
		fmt.Println(rows)
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	rows.Close()
}

func dbViimeisimmatUniikit(userID int) []string {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	sql := "select kuvaus from uniikitkahvit where user_id=$1 limit 5"
	rows, _ := conn.Query(context.Background(), sql, userID)

	var list []string

	for rows.Next() {
		var kahvi string

		err := rows.Scan(&kahvi)
		if err != nil {
			log.Println(err)
			os.Exit(1)
		}
		list = append(list, kahvi)
	}

	return list
}

func dbKupit(id int) int {
	conn, err := pgx.Connect(context.Background(), os.Getenv("PSQL_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Laske kupilliset ja palauta lukumäärä
	var kuppeja int
	sql := "select kupit from kuppilaskuri where user_id=$1"
	err = conn.QueryRow(context.Background(), sql, id).Scan(&kuppeja)
	if err == pgx.ErrNoRows {
		kuppeja = 0
	} else if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	return kuppeja
}
