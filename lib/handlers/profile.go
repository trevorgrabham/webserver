package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/profile"
	profiletemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/profile"
)

type profileData struct {
	profile.UserDetails
	Errors					[]string
}

// TODO: handle the security sanitization for the profile pic file upload

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }

	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	profiletemplate.Profile.Execute(w, profileData{*client, []string{}})
}

func HandleEditPic(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	profiletemplate.EditPic.Execute(w, profileData{*client, []string{}})
}

func HandleSavePic(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	err = profile.AddProfilePic(client, r)
	var (
		defaultTooLarge *profile.ErrFileTooLarge
		defaultUnsupported *profile.ErrUnsupportedFileFormat
		defaultNoFile *profile.ErrNoFile
	)
	if errors.As(err, &defaultTooLarge) || errors.As(err, &defaultUnsupported) || errors.As(err, &defaultNoFile) {
		profiletemplate.EditPic.Execute(w, profileData{*client, []string{err.Error()}})
		return
	}
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	profiletemplate.ShowPic.Execute(w, profileData{*client, []string{}})
}

func HandleEditName(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	profiletemplate.EditName.Execute(w, profileData{*client, []string{}})
}

func HandleSaveName(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { http.Error(w, "Unable to parse form", http.StatusBadRequest); return }

	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	res, ok := r.Form["name"]
	if !ok { profiletemplate.EditName.Execute(w, profileData{*client, []string{"No 'name' provided."}}); return }

	err = database.UpdateClient(&profile.UserDetails{ ID: userID, Name: res[0] })
	if err != nil { http.Error(w, "Unable to update 'name' on the server", http.StatusInternalServerError); return }
	client.Name = res[0]

	profiletemplate.ShowName.Execute(w, profileData{*client, []string{}})
}

func HandleEditEmail(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	profiletemplate.EditEmail.Execute(w, profileData{*client, []string{}})
}

func HandleSaveEmail(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("handleSaveEmail(): %v", err)) }

	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	res, ok := r.Form["email"]
	if !ok { profiletemplate.EditEmail.Execute(w, profileData{*client, []string{"No 'email' provided."}}); return }

	err = database.UpdateClient(&profile.UserDetails{ ID: userID, Email: res[0] })
	if err != nil { http.Error(w, "Unable to update 'email' on the server", http.StatusInternalServerError); return }
	client.Email = res[0]

	profiletemplate.ShowEmail.Execute(w, profileData{*client, []string{}})
}