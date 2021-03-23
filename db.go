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

	sql := "select kuvaus, max(aika) as aika from juonnit where user_id=$1 group by kuvaus order by aika desc limit 5"

	rows, _ := conn.Query(context.Background(), sql, userID)

	var list []string

	for rows.Next() {
		var kahvi string
		var asd string
		rows.Scan(&kahvi, &asd)
		log.Println(kahvi)
		list = append(list, kahvi)
	}

	log.Println(list)

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
