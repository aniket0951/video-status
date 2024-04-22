package helper

import (
	"io/ioutil"
	"mime/multipart"
	"path"
	"regexp"
	"strings"
)

func ValidateNumber(numb string) bool {
	re := regexp.MustCompile(`[^0-9]*1[34578][0-9]{9}[^0-9]*`)
	if len(numb) > 13 {
		return false
	}
	return re.MatchString(numb)
}

func CheckErr(err string) string {
	if strings.Contains(err, "oneof") {
		return "Invalid tag has been detected!"
	}

	return "Something Went's Wrong"
}

func LocalFileWrite(file multipart.File, localFilePath, extension string) (string, string, error) {
	tempFile, err := ioutil.TempFile(localFilePath, extension)

	if err != nil {
		return "", "", err
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		return "", "", err
	}

	tempFile.Write(fileBytes)
	defer file.Close()

	// FileKey, FilePath, Error
	filePath, fileKey := tempFile.Name(), path.Base(tempFile.Name())
	return fileKey, filePath, nil
}
