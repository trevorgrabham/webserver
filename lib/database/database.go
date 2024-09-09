package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/trevorgrabham/webserver/webserver/lib/chart"
	"github.com/trevorgrabham/webserver/webserver/lib/dashboard"
	"github.com/trevorgrabham/webserver/webserver/lib/profile"
	tagpkg "github.com/trevorgrabham/webserver/webserver/lib/tag"
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

func AddDay(activity dashboard.ActivityMetaData) (id int64, err error) {
	if activity.Description == "" || activity.Duration < 0 || activity.Day == "" {
		return -1, fmt.Errorf("AddDay(%v): Activities not initialized", activity)
	}
	// insert activity
	res, err := DB.Exec("INSERT INTO timer_data (user_id, day, duration, description) VALUES (?, ?, ?, ?)", activity.UserID, activity.Day, activity.Duration, activity.Description)
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
		tagID, err := res.LastInsertId()
		if err != nil {
			return id, fmt.Errorf("AddDay(%v): %v", activity, err)
		}
		activity.Tags[i].ID = tagID
	}
	return
}

// De-duplicates any identical tags on the same day, but not over different sessions
func GetDayData(userID int64, day string) (activities []dashboard.ActivityMetaData, err error) {
	if day == "" {
		return nil, fmt.Errorf("GetDayData(%s): Empty 'day'", day)
	}
	rows, err := DB.Query("SELECT timer_data.id AS id, duration, description, tag FROM timer_data JOIN activity_tag ON timer_data.id = activity_tag.activity_id WHERE day LIKE ? AND user_id = ?", day, userID)
	if err != nil {
		return nil, fmt.Errorf("GetDayData(%v, %v): %v", userID, day, err)
	}
	defer rows.Close()
	var (
		activityID, duration, previousId int64
		tag, description string
	)
	activities = make([]dashboard.ActivityMetaData, 0)
	for rows.Next() {
		err = rows.Scan(&activityID, &duration, &description, &tag)
		if err != nil {
			return nil, err
		}
		/* 
			We can have multiple tags per activityId, so we use previousId to see if this current row is refering to the same activityId as the previous row. If it is then we just need to grab the tag, otherwise we need to add a new activity
		*/
		if activityID != previousId {
			if DEBUG {
				fmt.Printf("New activity for user %d on %s\n", userID, day)
			}
			previousId = activityID
			activities = append(activities, dashboard.ActivityMetaData{
				ID: activityID,
				UserID: userID,
				Duration: duration,
				Description: description,
				Day: day,
				Tags: []tagpkg.TagMetaData{{ID: -1, Tag: tag}}})
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
		activities[len(activities)-1].Tags = append(activities[len(activities)-1].Tags, tagpkg.TagMetaData{ID: -1, Tag: tag, Count: 1})
	}
	if rows.Err() != nil {
		err = fmt.Errorf("GetDayData(%v): %v", day, rows.Err())
	}
	return
}

// De-duplicates any identical tags on the same day, over differing sessions
func GetCardData(userID int64, maxItems int64) (cards []dashboard.CardMetaData, err error) {
	if maxItems < 0 { return nil, fmt.Errorf("GetCardData(%d, %d): Bad value for 'maxItems'", userID, maxItems) }
	if userID < 1 { return nil, fmt.Errorf("GetCardData(%d, %d): Bad value for 'userID'", userID, maxItems) }

	var rows *sql.Rows
	if maxItems > 0 {
		rows, err = DB.Query("SELECT DISTINCT day FROM timer_data WHERE user_id = ? ORDER BY day DESC LIMIT ?", userID, maxItems)
	} else {
		rows, err = DB.Query("SELECT DISTINCT day FROM timer_data WHERE user_id = ? ORDER BY day DESC", userID)
	}
	if err != nil { return nil, fmt.Errorf("GetCardData(%d): %v", maxItems, err) }
	defer rows.Close()

	var day string
	cards = make([]dashboard.CardMetaData, 0)
	for rows.Next() {
		err = rows.Scan(&day)
		if err != nil { return nil, fmt.Errorf("GetCardData(%d): %v", maxItems, err) }

		dayActivities, err := GetDayData(userID, day)
		if err != nil { return nil, fmt.Errorf("GetCardData(%d): %v", maxItems, err) }

		var totalMins int64 
		tags := make(tagpkg.Tags, 0)
		for _, a := range dayActivities {
			totalMins += a.Duration
			// loop through so that we can de-duplicate any tags that are spread out over differing activities on the same day
			for _, t := range a.Tags {
				if !tags.Contains(t.Tag) {
					tags = append(tags, t)
				}
			}
		}
		cards = append(cards, dashboard.CardMetaData{Activities: dayActivities, Tags: tags, TotalMins: totalMins, Day: day})
	}
	if rows.Err() != nil { return nil, fmt.Errorf("GetCardData(%d, %d): %v", userID, maxItems, err) }
	return
}

func GetTagData(userID int64, offset int64) (tags tagpkg.Tags, err error) {
	rows, err := DB.Query("SELECT tag, SUM(1) AS 'count' FROM activity_tag WHERE activity_id IN (SELECT id from timer_data WHERE user_id = ?) GROUP BY tag ORDER BY count DESC, tag LIMIT 10 OFFSET ?", userID, offset)
	if err != nil { return nil, fmt.Errorf("GetTagData(%v, %v): %v", userID, offset, err) }
	defer rows.Close()

	var (
		tag string 
		count int64
		maxCount int64
	)
	for rows.Next() {
		err := rows.Scan(&tag, &count)
		if err != nil { return nil, fmt.Errorf("GetTagData(%d, %d): %v", userID, offset, err) }

		tags = append(tags, tagpkg.TagMetaData{ID: -1, Tag: tag, Count: count})
		if count > maxCount {
			maxCount = count
		}
	}
	if rows.Err() != nil { return nil, fmt.Errorf("GetTagData(%d, %d): %v", userID, offset, rows.Err()) }
	
	for i := range tags {
		tags[i].MaxCount = maxCount
	}

	return
}

func GetPreviousTags(userID int64) (tags []string, err error) {
	rows, err := DB.Query("SELECT DISTINCT tag FROM activity_tag WHERE activity_id IN (SELECT id FROM timer_data WHERE user_id = ?)", userID)
	if err != nil { return nil, fmt.Errorf("GetPreviousTags(%d): %v", userID, err) }
	defer rows.Close()

	var tag string 
	for rows.Next() {
		err := rows.Scan(&tag)
		if err != nil { return nil, fmt.Errorf("GetPreviousTags(%d): %v", userID, err) }
		tags = append(tags, tag)
	}
	if rows.Err() != nil { return nil, fmt.Errorf("GetPreviousTags(%d): %v", userID, rows.Err()) }
	return 
}

func GetPreviousActivities(userID int64) (activites []string, err error) {
	rows, err := DB.Query("SELECT DISTINCT description FROM timer_data WHERE user_id = ?", userID)
	if err != nil { return nil, fmt.Errorf("GetPreviousActivities(%d): %v", userID, err) }
	defer rows.Close()

	var description string 
	for rows.Next() {
		err := rows.Scan(&description)
		if err != nil { return nil, fmt.Errorf("GetPreviousActivities(%d): %v", userID, err) }
		activites = append(activites, description)
	}
	if rows.Err() !=  nil { return nil, fmt.Errorf("GetPreviousActivities(%d): %v", userID, rows.Err()) }
	return
}

func AddUserID() (id int64, err error) {
	var res sql.Result
	res, err = DB.Exec(`INSERT INTO user (name, email) VALUES (NULL, NULL)`)
	if err != nil { return -1, err }
	id, err = res.LastInsertId()
	return
}

func UpdateUser(details *profile.UserDetails) error {
	if details == nil { return fmt.Errorf("UpdatetUser(%v): No 'details' provided", details) }
	if details.ID < 1 { return fmt.Errorf("UpdateUser(%v): Bad value for 'ID'", details) }
	var err error
	if details.Name != "" {
		_, err = DB.Exec(`UPDATE user SET name = ? WHERE id = ?`, details.Name, details.ID)
	}
	if err != nil { return fmt.Errorf("UpdateUser(%v): %v", details, err) }
	if details.Email != "" {
		row := DB.QueryRow(`SELECT 1 FROM user WHERE email = ?`, details.Email)
		if err := row.Scan(); err != sql.ErrNoRows { return &profile.ErrEmailAlreadyExists{Message: fmt.Sprintf("Email %s is already registered", details.Email) }}

		_, err = DB.Exec(`UPDATE user SET email = ? WHERE id = ?`, details.Email, details.ID)
	}
	if err != nil { return fmt.Errorf("UpdateUser(%v): %v", details, err) }
	return nil
}

func GetUser(userID int64) (user *profile.UserDetails, err error) {
	if userID < 1 { return nil, fmt.Errorf("getUser(%d): bad value for 'userID'", userID) }

	rows, err := DB.Query(`SELECT name, email FROM user WHERE id = ?`, userID)
	if err != nil { return nil, fmt.Errorf("getUser(%d): %v", userID, err) }
	defer rows.Close()

	var name, email sql.NullString
	for rows.Next() {
		err := rows.Scan(&name, &email)
		if err != nil { return nil, fmt.Errorf("getUser(%d): %v", userID, err) }
	}
	return &profile.UserDetails{ 
		ID: 	userID, 
		Name: name.String, 
		Email: email.String, 
	}, nil
}

func LinkUsers(baseID int64, idToLink int64) error {
	if baseID < 1 { return fmt.Errorf("LinkUsers(%d, %d): %d is a bad value for 'baseId'", baseID, idToLink, baseID) }
	if idToLink < 1 { return fmt.Errorf("LinkUsers(%d, %d): %d is a bad value for 'idToLink'", baseID, idToLink, idToLink) }
	_, err := DB.Exec(`UPDATE timer_data SET user_id = ? WHERE user_id = ?`, baseID, idToLink)
	if err != nil { return fmt.Errorf("LinkUsers(%d, %d): %v", baseID, idToLink, err) }
	return nil
}

func GetStartEndData(userID int64) (start, end *time.Time, err error) {
		row, err := DB.Query(`SELECT MIN(day) AS start, MAX(day) AS end FROM timer_data WHERE user_id = ?`, userID)
		if err != nil { return nil, nil, fmt.Errorf("GetStartEndData(%d): %v", userID, err) }

		var s, e sql.Null[string]
		for row.Next() {
			err := row.Scan(&s, &e)
			if err != nil { return nil, nil, fmt.Errorf("GetStartEndData(%d): %v", userID, err) }
			if !s.Valid || !e.Valid { return nil, nil, nil}
		}
		if row.Err() != nil { return nil, nil, fmt.Errorf("GetStartEndData(%d): %v", userID, err) }

		{
			t, err := time.Parse(chart.DateMask, s.V)
			if err != nil { return nil, nil, fmt.Errorf("GetStartEndData(%d): %v", userID, err) }
			start = &t
		}
		{
			t, err := time.Parse(chart.DateMask, e.V)
			if err != nil { return nil, nil, fmt.Errorf("GetStartEndData(%d): %v", userID, err) }
			end = &t
		}
		
		return
}

func GetChartData(userID int64, start, end *time.Time) (res []chart.Data, err error) {
	if start == nil || end == nil { return nil, fmt.Errorf("GetChartData(%d, %s, %s): Cannot use Nil value for 'start' or 'end'") }

	rows, err := DB.Query(`SELECT day, duration FROM timer_data WHERE user_id = ? AND day >= ? AND day <= ?`, userID, start.Format(chart.DateMask), end.Format(chart.DateMask))
	if err != nil { return nil, fmt.Errorf("GetChartData(%d, %s, %s): %v", userID, start.Format(chart.DateMask), end.Format(chart.DateMask), err) }

	var (
		day string
		duration int64
	)
	for rows.Next() {
		if err := rows.Scan(&day, &duration); err != nil { return nil, fmt.Errorf("GetChartData(%d, %s, %s): %v", userID, start.Format(chart.DateMask), end.Format(chart.DateMask), err) }
		res = append(res, chart.Data{ Duration: float64(duration)/60.0, Day: day })
	}
	if rows.Err() != nil { return nil, fmt.Errorf("GetChartData(%d, %s, %s): %v", userID, start.Format(chart.DateMask), end.Format(chart.DateMask), rows.Err()) }
	return
}