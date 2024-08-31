package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

func CheckIDCookie(w http.ResponseWriter, r *http.Request) (userID int64, err error) {
	cookie, err := r.Cookie("client-id")
	if err == http.ErrNoCookie {
		userID, err = database.AddClientID()
		if err != nil { return -1, fmt.Errorf("CheckIDCookie(): %v", err) }
		cookie = &http.Cookie{
			Name: "client-id",
			Value: fmt.Sprint(userID),
			Path:"/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Expires: time.Now().AddDate(10, 0, 0),					 	// 10 years from now. For most browsers, this will default to 400 days from now
		}
	} else {
		userID, err = strconv.ParseInt(cookie.Value, 0, 64)
		if err != nil { return -1, fmt.Errorf("CheckIDCookie(): %v", err) }
		// Update the expiry on the client's id cookie whenever they re-visit the site
		cookie.Expires = time.Now().AddDate(10, 0, 0)
		cookie.Path = "/"
	}
	http.SetCookie(w, cookie)
	return userID, nil
}

func LinkToAccount(accountUserID int64, w http.ResponseWriter, r *http.Request) error {
	if accountUserID < 1 { return fmt.Errorf("LinkToAccount(%d): Bad value for 'accountUserID'", accountUserID) }
	oldUserID, err := CheckIDCookie(w, r)
	if err != nil { return fmt.Errorf("LinkToAccount(%d):%v", accountUserID, err) }
	err = database.LinkUsers(accountUserID, oldUserID)
	if err != nil { return fmt.Errorf("LinkToAccount(%d):%v", accountUserID, err) }
	return nil
}