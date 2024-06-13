package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type activity struct {
	id					int64
	day					string
	duration 		int64
	description	string
}

func addActivity(db *sql.DB, a activity) (int64, bool) {
	result, err := db.Exec("INSERT INTO timer_data (day, duration, description) VALUES (?, ?, ?)", a.day, a.duration, a.description)
	if err != nil {
		fmt.Printf("Inserting into db: %v", err)
		return -1, false
	}
	if id, err := result.LastInsertId(); err != nil {
		fmt.Printf("Retrieving insert id: %v", err)
		return -1, false
	} else {
		return id, true
	}
}

func addTag(db *sql.DB, tag string, aid int64) bool {
	result, err := db.Exec("INSERT INTO activity_tag (tag, activity_id) VALUES (?, ?)", tag, aid)
	if err != nil {
		fmt.Printf("Inserting into db: %v", err)
		return false
	}
	if _, err := result.LastInsertId(); err != nil {
		fmt.Printf("Retrieving insert id: %v", err)
		return false
	}
	return true
}

func main() {
	// open input file
	if len(os.Args) != 2 {
		log.Fatalf("Wrong number of arguments: want 1")
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalf("Opening file: %v", err)
	}
	scanner := bufio.NewScanner(file)

	// open db connection 
	cfg := mysql.Config{
		User: 	os.Getenv("DBUSER"),
		Passwd:	os.Getenv("DBPASS"),
		Net: 		"tcp",
		Addr:		"127.0.0.1:3306",
		DBName:	"programming_time_tracker",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("Opening mysql: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Pinging db: %v", err)
	}
	fmt.Println("Connected!")

	// read in data

	defer file.Close()
	// var i int64 = 3
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.Replace(line, "\"", "", -1)
		fields := strings.Split(line, ",")

		// trim the white space from the description and tags
		for i, _ := range fields[2:] {
			fields[i+2] = strings.ToLower(strings.TrimSpace(fields[i+2]))
		}
		// add the activities to the timer_data db
		dur, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			log.Fatalf("ParseInt: %v", err)
		}
		a := activity{
			day:					fields[0],
			duration:			dur,	
			description: 	fields[2],
		}
		id, ok := addActivity(db, a)
		if !ok {
			log.Fatalf("Error in addActivity(%v)", a)
		}

		// add the tags to the activity_tag db
		tags := fields[3:]
		for _, tag := range tags {
			if ok := addTag(db, tag, id); !ok {
				log.Fatalf("Error addTag(%s, %d)", tag, id)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("scanner.Scan(): %v", err)
	}

	// on success, delete the .csv file so that we aren't duplicating entries into the db
	
	// close file manually so that it can be removed
	file.Close()
	if err := os.Remove(os.Args[1]); err != nil {
		log.Fatalf("Remove(%s): %v", os.Args[1], err)
	}
}