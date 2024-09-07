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
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleprofile()")) }

	user, err := database.GetUser(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	profileTemplate := template.Must(template.New("profile").ParseFiles(templateutil.ParseFiles["profile"]...))
	profileTemplate.Execute(w, templateutil.NewProfileData(*user, nil))
}

func HandleEditPic(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleeditpic()")) }
	user, err := database.GetUser(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	editPic := template.Must(template.New("editpic").ParseFiles(templateutil.ParseFiles["editpic"]...))
	editPic.Execute(w, templateutil.NewProfileData(*user, nil))
}

func HandleSavePic(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlesavepic()")) }
	user, err := database.GetUser(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	savePic := template.Must(template.New("savepic").ParseFiles(templateutil.ParseFiles["savepic"]...))
	err = profile.AddProfilePic(user, r)
	var (
		defaultTooLarge *profile.ErrFileTooLarge
		defaultUnsupported *profile.ErrUnsupportedFileFormat
		defaultNoFile *profile.ErrNoFile
	)
	if errors.As(err, &defaultTooLarge) || errors.As(err, &defaultUnsupported) || errors.As(err, &defaultNoFile) {
		err = savePic.Execute(w, templateutil.NewProfileData(*user, []string{err.Error()}))
		if err != nil { http.Error(w, "Unable to generate template", http.StatusBadRequest)}
		return
	}
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	err = savePic.Execute(w, templateutil.NewProfileData(*user, nil))
	if err != nil { http.Error(w, "Unable to generate template", http.StatusBadRequest)}
}

func HandleEditName(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleeditname()")) }
	user, err := database.GetUser(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	editName := template.Must(template.New("editname").ParseFiles(templateutil.ParseFiles["editname"]...))
	editName.Execute(w, templateutil.NewProfileData(*user, nil))
}

func HandleSaveName(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { http.Error(w, "Unable to parse form", http.StatusBadRequest); return }

userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlesavename()")) }
	user, err := database.GetUser(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	saveName := template.Must(template.New("savename").ParseFiles(templateutil.ParseFiles["savename"]...))
	res, ok := r.Form["name"]
	if !ok { saveName.Execute(w, templateutil.NewProfileData(*user, []string{"No 'name' provided."})); return }

	err = database.UpdateUser(&profile.UserDetails{ ID: userID, Name: res[0] })
	if err != nil { http.Error(w, "Unable to update 'name' on the server", http.StatusInternalServerError); return }
	user.Name = res[0]

	saveName.Execute(w, templateutil.NewProfileData(*user, nil))
}

func HandleEditEmail(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleeditemail()")) }
	user, err := database.GetUser(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	editEmail := template.Must(template.New("editemail").ParseFiles(templateutil.ParseFiles["editemail"]...))
	editEmail.Execute(w, templateutil.NewProfileData(*user, nil))
}

func HandleSaveEmail(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("handleSaveEmail(): %v", err)) }

	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlesaveemail()")) }
	user, err := database.GetUser(userID)
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }

	saveEmail := template.Must(template.New("saveemail").ParseFiles(templateutil.ParseFiles["saveemail"]...))
	res, ok := r.Form["email"]
	if !ok { saveEmail.Execute(w, templateutil.NewProfileData(*user, []string{"No 'email' provided."})); return }

	err = database.UpdateUser(&profile.UserDetails{ ID: userID, Email: res[0] })
	if err != nil { http.Error(w, "Unable to update 'email' on the server", http.StatusInternalServerError); return }
	user.Email = res[0]

	saveEmail.Execute(w, templateutil.NewProfileData(*user, nil))
}