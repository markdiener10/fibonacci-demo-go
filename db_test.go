package main

import (
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	//"github.com/ory/dockertest/v3"
)

//var gpool dockertest.Pool
//User github.com/ory/dockertest to spin up our local pg docker file
//Also should look at: https://docs.docker.com/engine/api/
func testdbdockerteststart() {
	/*
		if gpool != nil {
			return
		}
		gpool, err := dockertest.NewPool("")
		if err != nil {
			log.Fatalf("Could not connect to docker: %s", err)
		}
	*/
}

func testdb() (error, *sql.DB) {

	//Generate the connection string
	Pgconnstr()

	//Need to make sure the table does not exist
	db, err := sql.Open("postgres", Sqlconn)
	if err != nil {
		return err, nil
	}
	return nil, db
}

func TestCanWeCreateTheMemoizeTable(t *testing.T) {

	//Make sure the table is GONE
	err, gdb := testdb()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer gdb.Close()

	sql := `DROP TABLE IF EXISTS public.memoized;`
	_, err = gdb.Exec(sql)
	if err != nil {
		t.Error(err.Error())
		return
	}

	//Normal entry point
	err = Pgsetup()
	if err != nil {
		t.Error(err.Error())
		return
	}

	//Need to check for the existence of the memoize table
	sql = `SELECT COUNT(table_name) AS count FROM information_schema.tables	WHERE table_schema='public' AND table_name='memoized';`
	rows, err := gdb.Query(sql)
	if err != nil {
		t.Error(err.Error())
	}
	defer rows.Close()

	var count int = -1

	for rows.Next() {
		err = rows.Scan(&count)
		if err != nil {
			t.Error(err.Error())
		}
		break
	}

	if count != 1 {
		t.Error("Table Memoized was not correctly creadted:", count)
	}

}

func TestCanWeFetchValues(t *testing.T) {

	//Make sure the table is FLUSHED
	err, gdb := testdb()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer gdb.Close()

	sql := `DELETE FROM public.memoized;`
	sql += `INSERT INTO public.memoized VALUES(1,5);`
	sql += `INSERT INTO public.memoized VALUES(3,6);`
	_, err = gdb.Exec(sql)
	if err != nil {
		t.Error(err.Error())
		return
	}

	//Need to load and then fetch different numbers of values
	err, gmap := PgFetch(10)
	if err != nil {
		t.Error(err.Error())
	}

	if _, ok := gmap[1]; !ok {
		t.Error("Key does not exist:1")
		return
	}
	if _, ok := gmap[3]; !ok {
		t.Error("Key does not exist:3")
		return
	}

	if gmap[1] != 5 {
		t.Error("Value did not stick 1")
		return
	}
	if gmap[3] != 6 {
		t.Error("Value did not stick 3")
		return
	}

}

func TestDoesTheCacheCodeActuallyCache(t *testing.T) {

	//Make sure the table is GONE
	err, gdb := testdb()
	if err != nil {
		t.Error(err.Error())
		return
	}
	defer gdb.Close()

	sql := `DELETE FROM public.memoized;`
	_, err = gdb.Exec(sql)
	if err != nil {
		t.Error(err.Error())
		return
	}

	var g = map[uint]uint64{1: 2, 2: 3, 3: 4}
	err = PgCache(g)
	if err != nil {
		t.Error(err.Error())
	}

	err, gmap := PgFetch(10)
	if err != nil {
		t.Error(err.Error())
	}

	if _, ok := gmap[1]; !ok {
		t.Error("Key does not exist:1")
		return
	}
	if _, ok := gmap[2]; !ok {
		t.Error("Key does not exist:2")
		return
	}
	if _, ok := gmap[3]; !ok {
		t.Error("Key does not exist:3")
		return
	}

	if gmap[1] != 2 {
		t.Error("Value did not stick 2")
	}
	if gmap[2] != 3 {
		t.Error("Value did not stick 3")
	}
	if gmap[3] != 4 {
		t.Error("Value did not stick 4")
	}

}

func TestCountAccuracy(t *testing.T) {

	gfib := Tfibonacci{}
	gfib.Init()
	gfib.Memoized(18)

	//Use the code we already have to generate intermediate
	PgCache(gfib.Cache)

	//Need to validate several retrievals
	err, count := PgCount(120)
	if err != nil {
		t.Error(err.Error())
	}
	if count != 12 {
		t.Error("Count does not equal 12 for 120 value:", count)
	}

	err, count = PgCount(1598)
	if err != nil {
		t.Error(err.Error())
	}
	if count != 18 {
		t.Error("Count does not equal 18 for 1597 value:", count)
	}

}

func TestCanWeClearTheCache(t *testing.T) {

	//Put some values in there
	var g = map[uint]uint64{2: 2, 4: 3, 6: 4}
	PgCache(g)

	//Fetch those values
	err, gmapa := PgFetch(10)
	if err != nil {
		t.Error(err.Error())
		return
	}

	if len(gmapa) < 3 {
		t.Error("MapA is not long enought:", len(gmapa))
		return
	}

	err = PgClear()
	if err != nil {
		t.Error(err.Error())
		return
	}

	_, gmapb := PgFetch(10)

	if len(gmapb) > 0 {
		t.Error(err.Error())
	}

}
