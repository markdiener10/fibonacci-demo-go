package main

import "log"

func main() {

	err := PgCheck()
	if err != nil {
		log.Fatal(err)		
	}
	if err = Pgsetup(); err != nil {
		log.Fatal(err)
	}
	Webserver()

}
