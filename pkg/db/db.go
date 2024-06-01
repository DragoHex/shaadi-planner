package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

// SQL queries
const (
	tableColumns = `
	CREATE TABLE IF NOT EXISTS invitees (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	name_on_card TEXT,
	no_of_members INTEGER,
	address TEXT,
	phone_no TEXT,
	events TEXT,
	tags TEXT,
	gifts TEXT,
	note TEXT
	);`
	insertInviteeSQL = `INSERT INTO invitees (name, name_on_card, no_of_members, address, phone_no, events, tags, gifts, note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

)

var (
	dbFile    = filepath.Join("artifacts", "db", "data.db")
	csvImportBackup = filepath.Join("artifacts", "csv", "import.csv")
	csvExportBackup = filepath.Join("artifacts", "csv", "export.csv")
)

type SQLiteDB struct {
	// change to db connector
	db *sql.DB
}

// Import imports data from csv file
func (s *SQLiteDB) Import(csvPath string) error {
	// compare the old and new import
	log.Printf("comparing the file to be imported %s, with the previously imported %s", csvPath, csvImportBackup)
	if sameImport(csvPath) {
		log.Printf("the content of %v, is same as the last import, so skipping the operation", csvPath)
		return nil
	}

	// read from the csv file
	records, err := readCSV(csvPath)
	if err != nil {
		return fmt.Errorf("error reading csv file %s: %v", csvPath, err)
	}

	// create db if not exist
	if _, err := os.Stat(dbFile); err != nil {
		log.Printf("db file %s does not exist", dbFile)
		if err := os.MkdirAll(filepath.Dir(dbFile), os.ModePerm); err != nil {
			return fmt.Errorf("error creating db %s: %v", dbFile, err)
		}
		log.Printf("db file created at %s", dbFile)
	}

	// add columns to the db
	log.Println("adding data to the db")
	s.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("error in opening the db %s", dbFile)
	}
	defer s.db.Close()
	if _, err := s.db.Exec(tableColumns); err != nil {
		return fmt.Errorf("table creation failed for the db %s", dbFile)
	}

	// write to the db
	stmt, err := s.db.Prepare(insertInviteeSQL)
	if err != nil {
		return fmt.Errorf("error perparing sql statement for insertion: %v", err)
	}
	defer stmt.Close()

	for _, record := range records {
		if len(record) < 9 {
			return fmt.Errorf("invalid record: %v", record)
		}
		name := record[1]
		nameOnCard := record[2]
		noOfMembers := record[3]
		phoneNo := record[4]
		address := record[5]
		events := record[6]
		tags := record[7]
		gifts := record[8]
		notes := record[9]

		_, err := stmt.Exec(name, nameOnCard, noOfMembers, phoneNo, address, events, tags, gifts, notes)
		if err != nil {
			return fmt.Errorf("error inserting a record: %v", record)
		}
	}

	// create backup for the new import
	// will be used to compare subsequent imports
	// to avoid redundant db operations
	_ = backupImport(csvPath, "import")

	return nil
}
