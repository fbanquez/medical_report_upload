package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

// DatabasePath is a constant containig the main sqlite file path
const DatabasePath string = "../persistent/codeImages"

func convertFileToBase64(filePath string) (encoded string, err error) {
	// Open file on disk.
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Problems opening image file.", err)
	}

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		fmt.Println("Problems converting image file into byte slice.", err)
	}

	// Encode as base64.
	encoded = base64.StdEncoding.EncodeToString(content)

	return
}

func checkTable() (existsTable bool, err error) {
	//
	db, err := sql.Open("sqlite3", DatabasePath)
	if err != nil {
		fmt.Println("Error opening database.", err)
	}
	defer db.Close()

	statement := "SELECT CAST(COUNT(*) AS BIT) FROM sqlite_master WHERE type = 'table' AND name = 'code_images'"
	stmt, err := db.Prepare(statement)
	if err != nil {
		fmt.Println("Problems creating the statement that validates the existence of local storage.", err)
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&existsTable)
	if err != nil {
		fmt.Println("Could not find the tables within local storage.", err)
	}

	return
}

func createTable() (err error) {
	//
	db, err := sql.Open("sqlite3", DatabasePath)
	if err != nil {
		fmt.Println("Error opening database.", err)
	}
	defer db.Close()

	statement := `CREATE TABLE IF NOT EXISTS code_images (
		id INTEGER,
		content TEXT,
		PRIMARY KEY (id))`

	stmt, err := db.Prepare(statement)
	if err != nil {
		fmt.Println("Problems creating the local storage table. ", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec()

	return
}

func storeImage(id int, content string) (err error) {
	//
	db, err := sql.Open("sqlite3", DatabasePath)
	if err != nil {
		fmt.Println("Error opening database.", err)
	}
	defer db.Close()

	statement := `INSERT INTO code_images (id, content) VALUES ($1, $2)`

	stmt, err := db.Prepare(statement)
	if err != nil {
		fmt.Println("Problems building the statement that insert image's data into local storage. ", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, content)
	if err != nil {
		fmt.Println("An issue occurred while storing images within local storage. ", err)
	}

	return
}

func main() {
	//
	idImg := flag.Int("i", 0, "image's Id")
	imgPath := flag.String("p", "../", "image's path")
	flag.Parse()

	code, err := convertFileToBase64(*imgPath)
	if err != nil {
		fmt.Println("Problem converting image to base64 code. ", err)
	}

	exists, err := checkTable()

	if !exists {
		if err = createTable(); err != nil {
			fmt.Println("Problem creating local table of code_images. ", err)
		}
	}

	if err = storeImage(*idImg, code); err != nil {
		fmt.Println("Error inserting coded image into local storage. ", err)
	}
}
