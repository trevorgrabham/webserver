package handlers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/trevorgrabham/webserver/webserver/lib/database"
	profiletemplate "github.com/trevorgrabham/webserver/webserver/lib/templates/profile"
	"github.com/trevorgrabham/webserver/webserver/lib/user"
)

type profileData struct {
	user.UserDetails
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
	
	err = r.ParseMultipartForm(10 << 30)
	if err == multipart.ErrMessageTooLarge { 
		profiletemplate.EditPic.Execute(w, profileData{*client, []string{"File too large. Files must be < 10MB."}})
		return 
	}
	if err != nil { http.Error(w, "Unable to parse uploaded file", http.StatusBadRequest); return }

	file, _, err := r.FormFile("pic")
	if err != nil { http.Error(w, "Unable to retrieve uploaded file", http.StatusBadRequest); return }
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil { http.Error(w, "Error reading uploaded file", http.StatusBadRequest); return }

	filetype := http.DetectContentType(buf)
	if filetype != "image/jpeg" && filetype != "image/jpg" && filetype != "image/png" {
		profiletemplate.EditPic.Execute(w, profileData{*client, []string{"Unsupported file type. Must be a .jpg, .jpeg, or .png file."}})
		return
	}
	filetype = strings.Split(filetype, "/")[1]

	// Somewhere round here would be a good place to sanitize the uploaded file

	// If there was already a custom profile pic, delete it before creating a new one for the user
	_, err = os.Stat(fmt.Sprintf("./static/imgs/user-%d.%s", userID, client.Ext))
	if err == nil {
		err := os.Remove(fmt.Sprintf("./static/imgs/user-%d.%s", userID, client.Ext))
		if err != nil { http.Error(w, "Unable to remove old version of the profile picture", http.StatusBadRequest); return }
	}

	err = database.UpdateClient(&user.UserDetails{ID: userID, Ext: filetype})
	if err != nil { http.Error(w, err.Error(), http.StatusBadRequest); return }
	client.Ext = filetype

	file.Seek(0, io.SeekStart)

	dst, err := os.Create(fmt.Sprintf("./static/imgs/user-%d.%s", userID, filetype))
	if err != nil { http.Error(w, "Unable to upload file to the server", http.StatusInternalServerError); return }
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil { http.Error(w, "Unable to save the file to the server", http.StatusInternalServerError); return }

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

	err = database.UpdateClient(&user.UserDetails{ ID: userID, Name: res[0] })
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

	err = database.UpdateClient(&user.UserDetails{ ID: userID, Email: res[0] })
	if err != nil { http.Error(w, "Unable to update 'email' on the server", http.StatusInternalServerError); return }
	client.Email = res[0]

	profiletemplate.ShowEmail.Execute(w, profileData{*client, []string{}})
}