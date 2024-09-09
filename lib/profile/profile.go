package profile

import (
	"fmt"
	"io"
	"io/fs"
	"mime/multipart"
	"net/http"
	"os"
)

const DefaultImgPath = `default.jpg`

type UserDetails struct {
	ID    int64
	Name  string
	Email string
}

type ErrEmailAlreadyExists struct { Message string }
func (e *ErrEmailAlreadyExists) Error() string { return e.Message }

func GetProfilePic(value interface{}) (path string) {
	userID, ok := value.(int64)
	if !ok { return "" }
	if userID < 1 { return "" }

	// Check if there is a ./static/imgs/userID folder present (indicates that there was a custom profile pic uploaded)
	files, err := os.ReadDir(fmt.Sprintf(`./static/imgs/%d`, userID))
	if os.IsNotExist(err) { return DefaultImgPath }
	if err != nil { return "" }

	// Make sure that there isn't more than one profile pic saved
	if len(files) > 1 { 
		var mostRecent fs.FileInfo
		for i := range files {
			file, err := os.Stat(files[i].Name())
			if err != nil { return "" }
			if file.ModTime().Before(mostRecent.ModTime()) {
				os.Remove(file.Name())
				continue
			}
			os.Remove(mostRecent.Name())
			mostRecent = file
		}
		return fmt.Sprintf(`%d/%s`, userID, mostRecent.Name())
	}

	return fmt.Sprintf(`%d/%s`, userID, files[0].Name())
}

// Calling function should check if err == ErrFileTooLarge || ErrUnsupportedFileFormat (under the profile package) so that they can alert the user in these cases
func AddProfilePic(user *UserDetails, r *http.Request) (errors []string) {
	if user == nil { panic(fmt.Errorf("AddProfilePic(%v): No user provided", user)) }
	
	// Parse form
	err := r.ParseMultipartForm(10 << 20) 		// 10MB
	if err == multipart.ErrMessageTooLarge { errors = append(errors, "File too large. Max size is 10MB") }
	if err != nil && err != multipart.ErrMessageTooLarge { panic(fmt.Errorf("AddProfilePic(%v): %v", user, err)) }

	// Grab the file
	file, fileHandle, err := r.FormFile("pic")
	if err != nil { return append(errors, "No file provided") }
	defer file.Close()

	// Check the file type
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil { panic(fmt.Errorf("AddProfilePic(%v): %v", user, err)) }

	filetype := http.DetectContentType(buf)
	if filetype != "image/jpeg" && filetype != "image/jpg" && filetype != "image/png" {
		return append(errors, fmt.Sprintf("Unsupported type %s", filetype))
	}
	
	// If there were user side errors, we have gathered all the error info that was relevant, we should return here
	if errors != nil { return errors }

	// See if there is already a directory for the user
	files, err := os.ReadDir(fmt.Sprintf(`./static/imgs/%d`, user.ID))
	if err != nil && !os.IsNotExist(err) { panic(fmt.Errorf("AddProfilePic(%v): %v", user, err)) }

	if os.IsNotExist(err) {
		// If there is no directory, create one and transfer to file
		err = os.Mkdir(fmt.Sprintf(`./static/imgs/%d`, user.ID), os.ModePerm)
		if err != nil { panic(fmt.Errorf("AddProfilePic(%v): %v", user, err)) }
	} else {
		// If there is a directory, delete all of the previous files in it
		for _, f := range files {
			err = os.Remove(fmt.Sprintf(`./static/imgs/%d/%s`, user.ID, f.Name()))
			if err != nil { panic(fmt.Errorf("AddProfilePic(%v): %v", user, err)) }
		}
	}

	// Create the new destination for the uploaded file
	destFile, err := os.Create(fmt.Sprintf(`./static/imgs/%d/%s`, user.ID, fileHandle.Filename))
	if err != nil { panic(fmt.Errorf("AddProfilePic(%v): %v", user, err)) }
	defer destFile.Close()

	// Before we transfer the uploaded file, need to reset it to the start of the file (from when we were checking the file type)
	file.Seek(0, io.SeekStart)
	_, err = io.Copy(destFile, file)
	if err != nil { panic(fmt.Errorf("AddProfilePic(%v): %v", user, err)) }
	return nil
}