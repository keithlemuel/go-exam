package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	_ "github.com/go-sql-driver/mysql"
)

type EnvVariable struct {
	Value string
	Key   string
}

type FileUpload struct {
	Filename string
	Size     int64
	MimeType string
	FilePath string
}

func setRoutes() {
	http.HandleFunc("/uploadfile", uploadFile)
	http.HandleFunc("/index", initIndex)
	http.ListenAndServe(":8000", nil)
}

func initIndex(writer http.ResponseWriter, request *http.Request) {
	key := "auth"
	envVar := EnvVariable{Value: getEnvVariable(key), Key: key}
	t, _ := template.ParseFiles("fileupload.html")
	t.Execute(writer, envVar)
}

// db conn
func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := ""
	dbName := "go_exam"
	dbHostName := "127.0.0.1"
	dbPort := "3306"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@tcp("+dbHostName+":"+dbPort+")/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	return db
}

// Save file to db
func insertFile(fileUpload FileUpload) bool {
	db := dbConn()
	insForm, err := db.Prepare("INSERT INTO Uploads(filename, size, mimeType, filePath) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	insForm.Exec(fileUpload.Filename, fileUpload.Size, fileUpload.MimeType, fileUpload.FilePath)

	defer db.Close()
	return true
}

func setEnvVariable(envVar EnvVariable) {
	os.Setenv(envVar.Key, envVar.Value)
}

func getEnvVariable(key string) string {
	return os.Getenv(key)
}

func returnStatusCodeErr(writer http.ResponseWriter, errorCode string) {
	writer.WriteHeader(http.StatusForbidden)
	switch errorCode {
	case "500":
		writer.Write([]byte("500 HTTP Internal Server Error!"))
	case "403":
		writer.Write([]byte("403 HTTP Bad Request!"))
	}
}

// Validate if image
func checkIfImage(mimeType string) bool {
	imageMimeTypes := []string{"image/apng", "image/bmp", "image/gif", "image/jpeg", "image/png", "image/svg+xml", "image/tiff", "image/webp"}
	for _, imageMimeType := range imageMimeTypes {
		if imageMimeType == mimeType {
			return true
		}
	}
	return false
}

// Process File Upload
func uploadFile(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("file is uploading...")

	request.Body = http.MaxBytesReader(writer, request.Body, 8<<20)
	file, handler, err := request.FormFile("file")
	if err != nil {
		fmt.Fprintf(writer, "413 HTTP Request Body Too Large! Max file sixe is 8 MB \n")
		log.Println(err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		returnStatusCodeErr(writer, "500")
		fmt.Println(err)
	}

	detectFile := mimetype.Detect(fileBytes)
	mime := detectFile.String()
	extension := detectFile.Extension()

	if request.FormValue("auth") != getEnvVariable("auth") || !checkIfImage(strings.ToLower(mime)) {
		returnStatusCodeErr(writer, "403")
		return
	}

	if err != nil {
		returnStatusCodeErr(writer, "500")
		fmt.Println("Error in retrieving the file: ")
		fmt.Println(err)
		return
	}
	defer file.Close()

	tempFile, err := ioutil.TempFile("files", "temp-*"+extension)
	if err != nil {
		returnStatusCodeErr(writer, "500")
		fmt.Println(err)
	}
	defer tempFile.Close()

	filePath := tempFile.Name()

	tempFile.Write(fileBytes)

	fileUpload := FileUpload{Filename: handler.Filename, Size: handler.Size, MimeType: mime, FilePath: filePath}
	result := insertFile(fileUpload)

	if !result {
		returnStatusCodeErr(writer, "500")
		return
	}

	fmt.Fprintf(writer, "File Successfully Uploaded\n")
}

// MAIN FUNCTION
func main() {
	fmt.Printf("Go Exam Starting.....")
	envVar := EnvVariable{Value: "123456778huidhsfksjfbk", Key: "auth"}
	setEnvVariable(envVar)
	setRoutes()
}
