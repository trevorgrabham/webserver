package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/trevorgrabham/webserver/webserver/html"
	"github.com/trevorgrabham/webserver/webserver/lib/database"
	"github.com/trevorgrabham/webserver/webserver/lib/profile"
)

// TODO: handle the security sanitization for the profile pic file upload

func HandleProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleprofile()")) }

	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }

	userProfile := html.ProfileData{
		Pic: html.ProfileFieldData{Value: user.ID },
		Name: html.ProfileFieldData{Value: user.Name },
		Email: html.ProfileFieldData{Value: user.Email },
	}

	profileTemplate := template.Must(template.New("profile").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["profile"]...))
	if err := profileTemplate.Execute(w, userProfile); err != nil { panic(err) }
}

func HandleEditPic(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleeditpic()")) }

	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }

	editData := html.ProfileFieldData{ Value: user.ID }

	editPic := template.Must(template.New("edit-pic").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["edit-pic"]...))
	if err := editPic.Execute(w, editData); err != nil { panic(err) }
}

func HandleSavePic(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlesavepic()")) }

	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }

	errors := profile.AddProfilePic(user, r)
	saveData := html.ProfileFieldData{ Value: user.ID, Errors: errors }

	if errors == nil {
		showPic := template.Must(template.New("show-pic").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["show-pic"]...))
		if err = showPic.Execute(w, saveData); err != nil { panic(err) }
		return
	}
	editPic := template.Must(template.New("edit-pic").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["edit-pic"]...))
	if err = editPic.Execute(w, saveData); err != nil { panic(err) }
}

func HandleEditName(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleeditname()")) }

	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }

	editData := html.ProfileFieldData{ Value: user.Name }

	editName := template.Must(template.New("edit-name").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["edit-name"]...))
	if err := editName.Execute(w, editData); err != nil { panic(err) }
}

func HandleSaveName(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic("Unable to parse form") }

	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlesavename()")) }
	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }
	
	name := r.Form.Get("name")
	var errors []string
	if name == "" { errors = append(errors, "No name provided") }
	if ok := checkInputForHTML(name); !ok { 
		errors = append(errors, "Unaccepted special character found in input")

		data := html.ProfileFieldData{ Value: user.Name, Errors: errors }

		editName := template.Must(template.New("edit-name").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["edit-name"]...))
		if err := editName.Execute(w, data); err != nil { panic(err) }
		return
	}

	saveData := html.ProfileFieldData{ Value: user.Name, Errors: errors }

	if errors == nil {
		err = database.UpdateUser(&profile.UserDetails{ ID: userID, Name: name })
		if err != nil { panic("Unable to update 'name' on the server") }

		saveData.Value = name

		showName := template.Must(template.New("show-name").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["show-name"]...))
		if err := showName.Execute(w, saveData); err != nil { panic(err) }
		return
	}
	editName := template.Must(template.New("edit-name").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["edit-name"]...))
	if err := editName.Execute(w, saveData); err != nil { panic(err) }
}

func HandleEditEmail(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handleeditemail()")) }

	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }

	editData := html.ProfileFieldData{ Value: user.Email }

	editEmail := template.Must(template.New("edit-email").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["edit-email"]...))
	if err := editEmail.Execute(w, editData); err != nil { panic(err) }
}

func HandleSaveEmail(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil { panic(fmt.Errorf("handleSaveEmail(): %v", err)) }

	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("unable to parse 'user-id' from handlesaveemail()")) }

	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }

	email := r.Form.Get("email")
	var errors []string
	if email == "" { errors = append(errors, "No email provided") }
	if ok = checkInputForHTML(email); !ok { 
		errors = append(errors, "Unaccepted special character found in input") 
	}

	// returns ErrEmailAlreadyExists or a standard error if there was an error on our side
	err = database.UpdateUser(&profile.UserDetails{ ID: userID, Email: email })
	if err != nil { 
		emailError, ok := err.(*profile.ErrEmailAlreadyExists) 
		if !ok { panic("Unable to update 'email' on the server") }

		errors = append(errors, emailError.Error())
	} else { user.Email = email }

	saveData := html.ProfileFieldData{ Value: user.Email, Errors: errors }

	if errors == nil {
		showEmail := template.Must(template.New("show-email").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["show-email"]...))

		if err := showEmail.Execute(w, saveData); err != nil { panic(err) }
		return
	}
	editEmail := template.Must(template.New("edit-email").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["edit-email"]...))

	if err := editEmail.Execute(w, saveData); err != nil { panic(err) }
}

func HandleLinkAccount(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(ContextKey("user-id")).(int64) 
	if !ok { panic(fmt.Errorf("handlelinkaccount(): unable to parse 'user-id'")) }
	if err := r.ParseForm(); err != nil { panic(err) }

	user, err := database.GetUser(userID)
	if err != nil { panic(err.Error()) }

	var errors []string
	accountEmail := r.Form.Get("link-to-account")
	if accountEmail == "" { errors = append(errors, "No account to link to provided") }

	if errors == nil && template.HTMLEscapeString(accountEmail) != accountEmail {
		errors = append(errors, "Not a valid email address. Special characters are not allowed")
	}

	otherID, err := database.GetUserIDFromEmail(accountEmail)
	e, ok := err.(*profile.ErrNoEmailExists) 
	if ok { errors = append(errors, e.Error())}
	if err != nil && !ok { panic(err) }

	if errors == nil {
		if err = database.LinkUsers(otherID, userID); err != nil { panic(err) }
		if err = UpdateCookie(otherID, w, r); err != nil { panic(err) }

		// grab the new user data, and send them back to their new profile page
		user, err = database.GetUser(otherID)
		if err != nil { panic(err.Error()) }
	}
	userProfile := html.ProfileData{
		Pic: html.ProfileFieldData{Value: user.ID },
		Name: html.ProfileFieldData{Value: user.Name },
		Email: html.ProfileFieldData{Value: user.Email },
		Link: errors}

	profileTemplate := template.Must(template.New("profile").Funcs(html.ProfileFuncMap).ParseFiles(html.IncludeFiles["profile"]...))
	if err := profileTemplate.Execute(w, userProfile); err != nil { panic(err) }
}

func checkInputForHTML(userString string) (ok bool) {
	return userString == template.HTMLEscapeString(userString)
}