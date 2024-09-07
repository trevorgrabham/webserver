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

type ErrNoFile struct {
	Message string
}

func (e *ErrNoFile) Error() string { return e.Message }

type ErrFileTooLarge struct {
	Message string
} 

func (e *ErrFileTooLarge) Error() string { return e.Message }

type ErrUnsupportedFileFormat struct {
	Message string
}

func (e *ErrUnsupportedFileFormat) Error() string { return e.Message }

func GetProfilePic(userID int64) (path string) {
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
func AddProfilePic(user *UserDetails, r *http.Request) error {
	if user == nil { return fmt.Errorf("AddProfilePic(%v): No user provided", user)}
	
	// Parse form
	err := r.ParseMultipartForm(10 << 20) 		// 10MB
	if err == multipart.ErrMessageTooLarge { return &ErrFileTooLarge{Message: "File too large. Max size is 10MB"} }
	if err != nil { return fmt.Errorf("AddProfilePic(%v): %v", user, err) }

	// Grab the file
	file, fileHandle, err := r.FormFile("pic")
	if err != nil { return &ErrNoFile{Message: "No file provided"} }
	defer file.Close()

	// Check the file type
	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil { return fmt.Errorf("AddProfilePic(%v): %v", user, err) }

	filetype := http.DetectContentType(buf)
	if filetype != "image/jpeg" && filetype != "image/jpg" && filetype != "image/png" {
		return &ErrUnsupportedFileFormat{Message: fmt.Sprintf("Unsupported type %s", filetype)}
	}

	// See if there is already a directory for the user
	files, err := os.ReadDir(fmt.Sprintf(`./static/imgs/%d`, user.ID))
	if err != nil && !os.IsNotExist(err) { return fmt.Errorf("AddProfilePic(%v): %v", user, err) }

	if os.IsNotExist(err) {
		// If there is no directory, create one and transfer to file
		err = os.Mkdir(fmt.Sprintf(`./static/imgs/%d`, user.ID), os.ModePerm)
		if err != nil { return fmt.Errorf("AddProfilePic(%v): %v", user, err) }
	} else {
		// If there is a directory, delete all of the previous files in it
		for _, f := range files {
			err = os.Remove(fmt.Sprintf(`./static/imgs/%d/%s`, user.ID, f.Name()))
			if err != nil { return fmt.Errorf("AddProfilePic(%v): %v", user, err) }
		}
	}

	// Create the new destination for the uploaded file
	destFile, err := os.Create(fmt.Sprintf(`./static/imgs/%d/%s`, user.ID, fileHandle.Filename))
	if err != nil { return fmt.Errorf("AddProfilePic(%v): %v", user, err) }
	defer destFile.Close()

	// Before we transfer the uploaded file, need to reset it to the start of the file (from when we were checking the file type)
	file.Seek(0, io.SeekStart)
	_, err = io.Copy(destFile, file)
	if err != nil { return fmt.Errorf("AddProfilePic(%v): %v", user, err)}
	return nil
}