package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/profile"
	"github.com/trevorgrabham/webserver/webserver/lib/templateutil"
)

// TODO: handle the security sanitization for the profile pic file upload

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }

	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	profileTemplate := template.Must(template.New("profile").ParseFiles(templateutil.ParseFiles["profile"]...))
	profileTemplate.Execute(w, templateutil.NewProfileData(*client, nil))
}

func HandleEditPic(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	editPic := template.Must(template.New("editpic").ParseFiles(templateutil.ParseFiles["editpic"]...))
	editPic.Execute(w, templateutil.NewProfileData(*client, nil))
}

func HandleSavePic(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	savePic := template.Must(template.New("savepic").ParseFiles(templateutil.ParseFiles["savepic"]...))
	err = profile.AddProfilePic(client, r)
	var (
		defaultTooLarge *profile.ErrFileTooLarge
		defaultUnsupported *profile.ErrUnsupportedFileFormat
		defaultNoFile *profile.ErrNoFile
	)
	if errors.As(err, &defaultTooLarge) || errors.As(err, &defaultUnsupported) || errors.As(err, &defaultNoFile) {
		err = savePic.Execute(w, templateutil.NewProfileData(*client, []string{err.Error()}))
		if err != nil { http.Error(w, "Unable to generate template", http.StatusBadRequest)}
		return
	}
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	err = savePic.Execute(w, templateutil.NewProfileData(*client, nil))
	if err != nil { http.Error(w, "Unable to generate template", http.StatusBadRequest)}
}

func HandleEditName(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	editName := template.Must(template.New("editname").ParseFiles(templateutil.ParseFiles["editname"]...))
	editName.Execute(w, templateutil.NewProfileData(*client, nil))
}

func HandleSaveName(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { http.Error(w, "Unable to parse form", http.StatusBadRequest); return }

	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	saveName := template.Must(template.New("savename").ParseFiles(templateutil.ParseFiles["savename"]...))
	res, ok := r.Form["name"]
	if !ok { saveName.Execute(w, templateutil.NewProfileData(*client, []string{"No 'name' provided."})); return }

	err = database.UpdateClient(&profile.UserDetails{ ID: userID, Name: res[0] })
	if err != nil { http.Error(w, "Unable to update 'name' on the server", http.StatusInternalServerError); return }
	client.Name = res[0]

	saveName.Execute(w, templateutil.NewProfileData(*client, nil))
}

func HandleEditEmail(w http.ResponseWriter, r *http.Request) {
	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	editEmail := template.Must(template.New("editemail").ParseFiles(templateutil.ParseFiles["editemail"]...))
	editEmail.Execute(w, templateutil.NewProfileData(*client, nil))
}

func HandleSaveEmail(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("handleSaveEmail(): %v", err)) }

	userID, err := CheckIDCookie(w, r)
	if err != nil { http.Error(w, "Unable to read 'client-id' from user agent", http.StatusBadRequest); return }
	client, err := database.GetClient(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	saveEmail := template.Must(template.New("saveemail").ParseFiles(templateutil.ParseFiles["saveemail"]...))
	res, ok := r.Form["email"]
	if !ok { saveEmail.Execute(w, templateutil.NewProfileData(*client, []string{"No 'email' provided."})); return }

	err = database.UpdateClient(&profile.UserDetails{ ID: userID, Email: res[0] })
	if err != nil { http.Error(w, "Unable to update 'email' on the server", http.StatusInternalServerError); return }
	client.Email = res[0]

	saveEmail.Execute(w, templateutil.NewProfileData(*client, nil))
}