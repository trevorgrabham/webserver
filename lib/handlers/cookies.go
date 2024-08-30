package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

func CheckIDCookie(w http.ResponseWriter, r *http.Request) (userID int64, err error) {
	if cookie, e := r.Cookie("client-id"); e == http.ErrNoCookie {
		userID, err = database.AddClientID()
		if err != nil { return -1, fmt.Errorf("CheckIDCookie(): %v", err) }
		cookie = &http.Cookie{
			Name: "client-id",
			Value: fmt.Sprint(userID),
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
		}
		http.SetCookie(w, cookie)
		return userID, nil
	} else if e != nil {
		return -1, fmt.Errorf("CheckIDCookie(): %v", e)
	} else {
		userID, err = strconv.ParseInt(cookie.Value, 0, 64)
		if err != nil { return -1, fmt.Errorf("CheckIDCookie(): %v", err) }
		return  
	}
}

func LinkToAccount(accountUserID int64, w http.ResponseWriter, r *http.Request) error {
	if accountUserID < 1 { return fmt.Errorf("LinkToAccount(%d): Bad value for 'accountUserID'", accountUserID) }
	oldUserID, err := CheckIDCookie(w, r)
	if err != nil { return fmt.Errorf("LinkToAccount(%d):%v", accountUserID, err) }
	err = database.LinkUsers(accountUserID, oldUserID)
	if err != nil { return fmt.Errorf("LinkToAccount(%d):%v", accountUserID, err) }
	return nil
}