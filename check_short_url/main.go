package main

import (
	"database/sql"
	"fmt"
	"time"
	"log"
	"net/http"
	_ "github.com/lib/pq"
	"sync"
	"context"
)	

var ctx context.Context
var db, db_err = sql.Open("postgres", "user=root password=root dbname=short_url sslmode=disable")
const INTERVAL = 1 // minutes

func main() {
	checkConnetionBd()
	defer db.Close()
	
	queryGetUrls := fmt.Sprintf("SELECT id, short_url FROM url_url WHERE ((last_check < NOW() - INTERVAL '%d minutes') or (last_check is NULL));", INTERVAL) 
	queryCountUrls := fmt.Sprintf("SELECT count(*) FROM url_url WHERE ((last_check < NOW() - INTERVAL '%d minutes') or (last_check is NULL));", INTERVAL)
	queryGetOlderCheckUrl := "SELECT last_check FROM url_url ORDER BY last_check ASC"
	var countUrls int

	for {
		err := db.QueryRow(queryCountUrls).Scan(&countUrls)   
		if err != nil {
			log.Fatal(err)
		}

		if countUrls > 0 {
			createGroupsToCheckUrls(countUrls, queryGetUrls)
		} else {
			waitToCheckUrls(queryGetOlderCheckUrl)
		}
	}
}

func createGroupsToCheckUrls(countUrls int, queryGetUrls string) {
	numGroups := 10
	if countUrls < 10 {
		numGroups = countUrls
	}

	fmt.Printf("Num of groups: %v \n", numGroups)

	rows, err := db.Query(queryGetUrls)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for i := 0; i < numGroups; i++ {
		wg.Add(1)
		go func(jobID int) {
			checkUrls(rows, jobID)
			defer wg.Done()
		}(i)
	}

	wg.Wait()
}

func waitToCheckUrls(queryGetOlderCheckUrl string) {
	var olderTime time.Time
	err := db.QueryRow(queryGetOlderCheckUrl).Scan(&olderTime)
	if err != nil {
		log.Fatal(err)
	}

	timeWait := (olderTime.Add(INTERVAL * time.Minute)).Sub(time.Now())

	fmt.Printf("Sleeping for %v \n", timeWait)
	fmt.Printf("Older Time %v \n", olderTime)
	time.Sleep(timeWait)
}

func checkConnetionBd() {
	if db_err != nil {
		log.Fatal(db_err)
	}

	db_err = db.Ping()
	if db_err != nil {
		log.Fatal(db_err)
	}
}

func checkUrls(urls *sql.Rows, jobID int) {
	for urls.Next() {
		var url string
		var id int
		if err := urls.Scan(&id, &url); err != nil {
				log.Fatal(err)
		}
		statusCode := makeRequestAndUpdate(id, url)
		fmt.Printf("ID: %d - %s - Job ID %d - Status %d\n", id, url, jobID, statusCode)
	}
}

func makeRequestAndUpdate(id int, url string) int{
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

	return statusCode
}
