package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/* NOTE Depending on style preference, some golang teams like fixed routing
http.HandleFunc("/", func())
http.HandleFunc("/CLEAR", clearfunc())
http.HandleFunc("/FETCH", fetchfunc())
http.HandleFunc("/COUNT", countfunc())
log.Fatal(http.ListenAndServe(":8081", nil))
*/

func Webserver() {

	srv := &http.Server{
		Addr:           ":8081",
		Handler:        http.HandlerFunc(webhandler),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		log.Println("Starting Server")
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal("Server Returned:", err.Error())
		}
	}()

	intchan := make(chan os.Signal, 1)
	signal.Notify(intchan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-intchan

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	srv.Shutdown(ctx)
	log.Println("Shutting down...")
	os.Exit(0)
}

//Bug Report
//#93
//12200160415121876738
//12200160415121876738

//#94
//19740274219868223167
//1293530146158671551 //Running into max value

//Internationalization anyone?
const badcmd = "<br>Valid Commands are:<br>http://localhost:8081/FETCH/XX (ordinal) - Produce the fibonacci value given an ordinal.<br>Maximum Value is 90 Ordinal.<br>http://localhost:8081/COUNT/XX (fibonacci) - Number of intermediates required to produce given fibonacci value.<br>Maximum Value is 12200160415121876738 Fibonacci.<br>http://localhost:8081/CLEAR - Clear Data Store"

func webhandler(w http.ResponseWriter, r *http.Request) {

	parms := strings.Split(r.URL.Path, "/")
	var val uint64 = 0
	cmd := "NONE"
	sval := ""

	if len(parms) < 0 {
		w.Write([]byte(badcmd))
		return
	}

	for _, v := range parms {
		sval = strings.ToUpper(v)
		switch sval {
		case "":
			continue
		case "FETCH":
			fallthrough
		case "COUNT":
			fallthrough
		case "CLEAR":
			cmd = sval
			continue
		default:
			//Try to convert to a number
			u64, err := strconv.ParseUint(sval, 10, 64)
			if err != nil {
				continue
			}
			if u64 > val {
				val = u64
			}
		}
	}

	if cmd == "CLEAR" {
		w.Write([]byte("CLEAR Command received "))
		err := PgClear()
		if err != nil {
			w.Write([]byte(fmt.Sprint("<br>ERROR: Data Store Error Detected:", err.Error())))
		}
		return
	}

	switch cmd {
	case "FETCH":

		//Put an arbitrary clamp on the maximum value
		if val > 90 {
			val = 90
		}

		//Fetch the cache first
		err, gafibo := PgFetch(val)
		if err != nil {
			w.Write([]byte(fmt.Sprint("<br>ERROR: Data Store Error Detected:", err.Error())))
			gafibo = make(map[uint]uint64)
		}

		//Pre-load the fibonacci generator
		gfib := Tfibonacci{}
		gfib.Init()
		for key, val := range gafibo {
			gfib.Cache[key] = val
		}

		//Run the pre-loaded memoized generator
		fibval := gfib.Memoized(uint(val))

		w.Write([]byte(fmt.Sprint("<br>FETCH Fibonacci:", val, " ordinal produces ", fibval, " Fibonacci value\r\n")))

		err = PgCache(gfib.Cache)
		if err != nil {
			w.Write([]byte(fmt.Sprint("<br>ERROR:Data Store Error Detected:", err.Error())))
		}

	case "COUNT":

		gfib := Tfibonacci{}
		gfib.Init()
		idx := uint(0)

		//Fibonacci 90 => Big value
		if val > 2880067194370816120 {
			val = 2880067194370816120
		}

		//We need to check the cache for values first
		err, count := PgCount(val)
		if err != nil {
			w.Write([]byte(fmt.Sprint("<br>ERROR:Data Store Error Detected:", err.Error())))
		}
		if count > 0 {
			w.Write([]byte(fmt.Sprint("<br>COUNT Command received:", val, " yields ", count, " memoized results")))
			return
		}

		//Need to use the fibonacci generator to iterate to the desired value
		for {
			if gfib.Memoized(idx) < val {
				idx++
				continue
			}
			break
		}
		w.Write([]byte(fmt.Sprint("<br>ERROR:COUNT Command received:", val, " yields ", idx, " memoized results")))

		//Cache the fibonacci intermediate values
		err = PgCache(gfib.Cache)
		if err != nil {
			w.Write([]byte(fmt.Sprint("<br>ERROR:Data Store Error Detected:", err.Error())))
		}
	default:
		w.Write([]byte(badcmd))
	}
}
