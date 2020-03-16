package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"sync"
	"time"
)

const INTERVAL = 1 // minutes
const GROUPWORKERS = 10

var ctx context.Context
var db, db_err = sql.Open("postgres", "user=postgres password=postgres dbname=short_url sslmode=disable")

type Url struct {
	Id  int
	Url string
	Status int
}

func main() {
	if db_err != nil {
		log.Fatal(db_err)
	}
	defer db.Close()

	queryGetUrls := fmt.Sprintf("SELECT id, short_url FROM url_url WHERE ((last_check < NOW() - INTERVAL '%d minutes') or (last_check is NULL));", INTERVAL)

	for {
		rows, err := db.Query(queryGetUrls)
		if err != nil {
			log.Fatal(err)
		}
		defer rows.Close()

		var urls []Url
		var id int
		var short_url string
		for rows.Next() {
			err := rows.Scan(&id, &short_url)
			if err != nil {
				log.Fatal(err)
			}

			urls = append(urls, Url{Id: id, Url: short_url})
		}
		if err := rows.Err(); err != nil {
			log.Fatal(err)
		}
		if len(urls) > 0 {
			createGroupsToCheckUrls(urls)
		} else {
			fmt.Println("Sleeping for 1 minute")
			time.Sleep(time.Minute)
		}
	}
}

func createGroupsToCheckUrls(urls []Url) {
	var wg sync.WaitGroup
    job := make(chan Url)
    result := make(chan bool)

	for w := 1; w <= GROUPWORKERS; w++ {
        go checkUrl(w, job, result)
	}

	for _, url := range urls {
		job <- url
	}

	close(job)

	for a := 1; a <= len(urls); a++ {
        <-result
	}

	wg.Wait()
}

func checkUrl(id int, job <-chan Url, result chan<- bool) {
	for url := range job {
	    fmt.Println("worker", id, "started ", url.Url)
		makeRequestAndUpdate(url.Id, url.Url)
		fmt.Println("worker", id, "finished ", url.Url)
        result <- true
	}
}

func makeRequestAndUpdate(id int, url string) {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	resp, err := client.Do(req)
	timeNow := time.Now()
	queryUpdateStatus := "UPDATE url_url SET last_check = $2, last_check_status = $3 WHERE id = $1"

	var statusCode int
	if err != nil {
		statusCode = 999
	} else {
		statusCode = resp.StatusCode
	}

	_, err = db.Exec(queryUpdateStatus, id, timeNow, statusCode)
	if err != nil {
		log.Fatal(err)
	}
}
