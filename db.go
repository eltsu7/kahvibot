package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx"
)

func dbViim(userID int) []string {
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
