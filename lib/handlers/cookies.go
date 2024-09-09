package handlers

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
)

var sc *securecookie.SecureCookie
type ContextKey string

func init() {
	hashString := os.Getenv("SCHASH")
	if hashString == "" { panic(fmt.Errorf("no 'SCHASH' env variable set")) }
	blockString := os.Getenv("SCBLOCK")
	if blockString == "" { panic(fmt.Errorf("no 'SCBLOCK' env variable set")) }
	hash, err := base64.StdEncoding.DecodeString(hashString)
	if err != nil { panic(err) }
	block, err := base64.StdEncoding.DecodeString(blockString)
	if err != nil { panic(err) }
	sc = securecookie.New(hash, block)
}

func SetCookieContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var userID int64
		cookie, err := r.Cookie("user-id")
		if err == http.ErrNoCookie {
			userID, err = database.AddUserID()
			if err != nil { panic(err) }
			encoded, err := sc.Encode("user-id", userID)
			if err != nil { panic(err) }
			cookie = &http.Cookie{
				Name: "user-id",
				Value: encoded,
				Path: "/",
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Secure: false,
				Expires: time.Now().AddDate(2, 0, 0)}
			http.SetCookie(w, cookie)
		} else if err = sc.Decode("user-id", cookie.Value, &userID); err != nil { panic(err) }
		ctx := context.WithValue(r.Context(), ContextKey("user-id"), userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LinkToAccount(accountUserID int64, w http.ResponseWriter, r *http.Request) error {
	oldUserID, ok := r.Context().Value(ContextKey("user-id")).(int64)
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' cookie")) }
	if accountUserID < 1 { return fmt.Errorf("LinkToAccount(%d): Bad value for 'accountUserID'", accountUserID) }
	err := database.LinkUsers(accountUserID, oldUserID)
	if err != nil { return fmt.Errorf("LinkToAccount(%d):%v", accountUserID, err) }
	return nil
}