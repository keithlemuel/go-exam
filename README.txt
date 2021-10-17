SETUP

1. Install go 
2. Update environment variable $GOPATH to C:\Go\go-works
3. In go-works folder, make sure to create bin, pkg and src folders
4. Create GoFileUpload folder
5. Create files folder under the GoFileUpload that you created
6. Exec these set of commands: 
    - go get -u github.com/go-sql-driver/mysql
        - This driver is need in order to connect to database and to handle executing queries and statements
    - go get github.com/gabriel-vasile/mimetype
        - This driver would help determine the extension and mime type of an file
7. I used xampp and workbench for the db (127.0.0.1:3306)
8. Crete Table using these commands:
    - CREATE SCHEMA go_exam
    - CREATE TABLE go_exam.uploads (
        id int(6) NOT NULL AUTO_INCREMENT,
        filename varchar(30) NOT NULL,
        size integer NOT NULL,
        mimeType varchar(30) NOT NULL,
        filePath varchar(30) NOT NULL,
        PRIMARY KEY (id)
    ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=latin1;


RUNNING THE PROJECT
1. Run this command in terminal: go run test.go and you should the "Go Exam Starting..." in the terminal
2. Go to this url http://localhost:8000/index
    - This will redirect to the upload page