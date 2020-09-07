package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	_ "os"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

const (
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

var Sqlconn string

//A function to allow docker-compose to override the localhost value and map in a different value
func pghost() string {
	pghost := "localhost"
	if os.Getenv("PGHOST") != "" {
		pghost = os.Getenv("PGHOST")
	}
	return pghost
}

func Pgconnstr() {
	if Sqlconn != "" {
		return
	}
	//Default value (undefined normally in dev mode)
	Sqlconn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", pghost(), port, user, password, dbname)
}

//We need a way to check to see if we can make a connection on 5432
func PgCheck() error {

	//When docker-compose comes up, we need a way to see if we can make a connection
	gconn := fmt.Sprintf("%s:%d", pghost(), port)

	for cnt := 0; cnt < 10; cnt++ {
		time.Sleep(500 * time.Millisecond)
		log.Println("PgCheck Tik:", cnt)
		conn, err := net.Dial("tcp", gconn)
		if err != nil {
			if strings.Contains(err.Error(), "connection refused") == true {
				continue
			}
			if strings.Contains(err.Error(), "EOF") == true {
				continue
			}
			if strings.Contains(err.Error(), "host") == true {
				log.Println(err.Error())
				continue
			}

			return err
		}
		defer conn.Close()
		break
		// handle error
	}

	return nil

}

//Setup our database stuff for this pass
func Pgsetup() error {
	//Make sure the table is created and ready for the test

	Pgconnstr()

	db, err := sql.Open("postgres", Sqlconn)
	if err != nil {
		return err
	}
	defer db.Close()
	sql := `DROP TABLE IF EXISTS public.memoized;`
	sql += `CREATE TABLE public.memoized(idx bigint PRIMARY KEY,fibo integer);`
	_, err = db.Exec(sql)
	return err
}

//Fetch our cached fibonacci values
func PgFetch(val uint64) (error, map[uint]uint64) {

	Pgconnstr()

	db, err := sql.Open("postgres", Sqlconn)
	if err != nil {
		return err, nil
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT idx,fibo FROM public.memoized WHERE idx <= %d ORDER BY IDX ASC;", val)
	rows, err := db.Query(sql)
	if err != nil {
		return err, nil
	}
	defer rows.Close()

	var g = make(map[uint]uint64)
	var idx uint
	var fibo uint64

	for rows.Next() {
		err = rows.Scan(&idx, &fibo)
		if err != nil {
			return err, nil
		}
		g[idx] = fibo
	}
	return nil, g
}

func PgCache(g map[uint]uint64) error {

	if len(g) == 0 {
		return nil
	}

	Pgconnstr()

	db, err := sql.Open("postgres", Sqlconn)
	if err != nil {
		return err
	}
	defer db.Close()

	sql := ""
	tsql := "INSERT INTO public.memoized VALUES(%d,%d) "
	tsql += "ON CONFLICT(idx) DO "
	tsql += "UPDATE SET fibo = %d WHERE public.memoized.idx = %d;"

	//Build up our INSERT/UPDATE multi query SQL
	for idx, fibo := range g {
		sql += fmt.Sprintf(tsql, idx, fibo, fibo, idx)
	}
	rows, err := db.Query(sql)
	if err != nil {
		return err
	}
	defer rows.Close()
	return nil
}

//I know, the int vs uint value store overflow for the int return value below
func PgCount(val uint64) (error, int) {

	Pgconnstr()

	db, err := sql.Open("postgres", Sqlconn)
	if err != nil {
		return err, 0
	}
	defer db.Close()

	sql := fmt.Sprintf("SELECT COUNT(idx) AS count FROM public.memoized WHERE fibo <= %d;", val)
	rows, err := db.Query(sql)
	if err != nil {
		return err, 0
	}
	defer rows.Close()
	var count int = -1
	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			return err, 0
		}
		break
	}
	return nil, count
}

func PgClear() error {

	Pgconnstr()

	db, err := sql.Open("postgres", Sqlconn)
	if err != nil {
		return err
	}
	defer db.Close()
	sql := `DELETE FROM public.memoized;`
	_, err = db.Exec(sql)
	return err
}
