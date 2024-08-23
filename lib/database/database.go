package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/go-sql-driver/mysql"
	"github.com/trevorgrabham/webserver/webserver/lib/util"
)

const DEBUG = false

var DB *sql.DB

var cfg = mysql.Config{
	User:		os.Getenv("DBUSER"),
	Passwd:	os.Getenv("DBPASS"),
	Net:		"tcp",
	Addr:		"127.0.0.1:3306",
	DBName:	"testing_db",
}

func init() {
	var err error
	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	if err := DB.Ping(); err != nil {
		log.Fatal(err)
	}
}

func AddDay(activity util.ActivityMetaData) (id int64, err error) {
	if activity.Description == "" || activity.Duration < 0 || activity.Day == "" {
		return -1, fmt.Errorf("AddDay(%v): Activities not initialized", activity)
	}
	// insert activity
	res, err := DB.Exec("INSERT INTO timer_data (day, duration, description) VALUES (?, ?, ?)", activity.Day, activity.Duration, activity.Description)
	if err != nil {
		return -1, fmt.Errorf("AddDay(%v): %v", activity, err)
	}
	id, err = res.LastInsertId()
	if err != nil {
		return -1, fmt.Errorf("AddDay(%v): %v", activity, err)
	}
	if DEBUG {
		fmt.Printf("Insert into DB for day %v success!\nID for inserted row into timer_data is %d\n", activity.Day, id)
	}
	// insert tags
	for i := range activity.Tags {
		res, err := DB.Exec("INSERT INTO activity_tag (tag, activity_id) VALUES (?, ?)", activity.Tags[i].Tag, id)
		if err != nil {
			return id, fmt.Errorf("AddDay(%v): %v", activity, err)
		}
		tagId, err := res.LastInsertId()
		if err != nil {
			return id, fmt.Errorf("AddDay(%v): %v", activity, err)
		}
		activity.Tags[i].Id = tagId
	}
	return
}

// De-duplicates any identical tags on the same day, but not over different sessions
func GetDayData(day string) (activities []util.ActivityMetaData, err error) {
	if day == "" {
		return nil, fmt.Errorf("GetDayData(%s): Empty 'day'", day)
	}
	rows, err := DB.Query("SELECT timer_data.id AS id, duration, description, tag FROM timer_data JOIN activity_tag ON timer_data.id = activity_tag.activity_id WHERE day LIKE ?", day)
	if err != nil {
		return nil, fmt.Errorf("GetDayData(%v): %v", day, err)
	}
	defer rows.Close()
	var (
		activityId, duration, previousId int64
		tag, description string
	)
	activities = make([]util.ActivityMetaData, 0)
	for rows.Next() {
		err = rows.Scan(&activityId, &duration, &description, &tag)
		if err != nil {
			return nil, err
		}
		/* 
			We can have multiple tags per activityId, so we use previousId to see if this current row is refering to the same activityId as the previous row. If it is then we just need to grab the tag, otherwise we need to add a new activity
		*/
		if activityId != previousId {
			if DEBUG {
				fmt.Printf("New activity for %s\n", day)
			}
			previousId = activityId
			activities = append(activities, util.ActivityMetaData{
				Id: activityId,
				Duration: duration,
				Description: description,
				Day: day,
				Tags: []util.TagMetaData{{Id: -1, Tag: tag}}})
			continue
		}
		if activities[len(activities)-1].Tags.Contains(tag) {
			if DEBUG {
				fmt.Printf("Same activity for day %s, but the tag %s is already present, so skipping it\n", day, tag)
			}
			continue
		}
		if DEBUG {
			fmt.Printf("Same activity for day %s, adding new tag %s\n", day, tag)
		}
		activities[len(activities)-1].Tags = append(activities[len(activities)-1].Tags, util.TagMetaData{Id: -1, Tag: tag})
	}
	if rows.Err() != nil {
		err = fmt.Errorf("GetDayData(%v): %v", day, rows.Err())
	}
	return
}

// De-duplicates any identical tags on the same day, over differing sessions
func GetCardData(maxItems int64) (cards []util.CardMetaData, err error) {
	if maxItems < 0 {
		return nil, fmt.Errorf("GetCardData(%d): Bad value for 'maxItems'", maxItems)
	}
	var rows *sql.Rows
	if maxItems > 0 {
		rows, err = DB.Query("SELECT DISTINCT day FROM timer_data ORDER BY day DESC LIMIT ?", maxItems)
	} else {
		rows, err = DB.Query("SELECT DISTINCT day FROM timer_data ORDER BY day DESC")
	}
	if err != nil {
		return nil, fmt.Errorf("GetCardData(%d): %v", maxItems, err)
	}
	defer rows.Close()

	var day string
	cards = make([]util.CardMetaData, 0)
	for rows.Next() {
		err = rows.Scan(&day)
		if err != nil {
			return nil, fmt.Errorf("GetCardData(%d): %v", maxItems, err)
		}
		dayActivities, err := GetDayData(day)
		if err != nil {
			return nil, fmt.Errorf("GetCardData(%d): %v", maxItems, err)
		}
		var totalMins int64 
		tags := make(util.Tags, 0)
		for _, a := range dayActivities {
			totalMins += a.Duration
			// loop through so that we can de-duplicate any tags that are spread out over differing activities on the same day
			for _, t := range a.Tags {
				if !tags.Contains(t.Tag) {
					tags = append(tags, t)
				}
			}
		}
		cards = append(cards, util.CardMetaData{Activities: dayActivities, Tags: tags, TotalMins: totalMins, Day: day})
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("GetCardData(%d): %v", maxItems, err)
	}
	return
}