package db

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/DragoHex/shaadiPlanner/pkg/invitee"
	_ "github.com/mattn/go-sqlite3"
)

// SQL queries
const (
	tableColumns = `
	CREATE TABLE IF NOT EXISTS invitees (
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT NOT NULL,
	name_on_card TEXT,
	num_of_members INTEGER,
	address TEXT,
	phone_num TEXT,
	events TEXT,
	tags TEXT,
	gifts TEXT,
	note TEXT
	);`
	insertInviteeSQL = `INSERT INTO invitees (name, name_on_card, num_of_members, address, phone_num, events, tags, gifts, note) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	selectAllSQL     = `SELECT * FROM invitees`
)

var (
	dbFile          = filepath.Join("artifacts", "db", "data.db")
	csvImportBackup = filepath.Join("artifacts", "csv", "import.csv")
	csvExportBackup = filepath.Join("artifacts", "csv", "export.csv")
)

type SQLiteDB struct {
	// TODO: change to db connector
	// TODO: Add an Column interface that has queries as object
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

func (s *SQLiteDB) Export(csvPath string) error {
	var err error
	// when db doesn't exist
	if _, err = os.Stat(dbFile); err != nil {
		log.Printf("export: db file %s does not exist", dbFile)
		return fmt.Errorf("db %s doesn't exist, exiting the export", dbFile)
	}

	// open a db connection
	s.db, err = sql.Open("sqlite3", dbFile)
	if err != nil {
		return fmt.Errorf("error i nopening the db %s: %v", dbFile, err)
	}
	defer s.db.Close()

	// query the rows
	rows, err := s.db.Query(selectAllSQL)
	if err != nil {
		return fmt.Errorf("error fetching data from database: %v", err)
	}
	defer rows.Close()

	// check if the csv exist, create if not
	if _, err := os.Stat(csvPath); err != nil {
		log.Printf("csv file %s for the data to export to does not exist", csvPath)
		log.Printf("creating csv file %s to write data to", csvPath)
		dirPath := filepath.Dir(csvPath)
		if dirPath != "." {
			if err := os.MkdirAll(filepath.Dir(csvPath), os.ModePerm); err != nil {
				return fmt.Errorf("error creating csv %s: %v", csvPath, err)
			}
			log.Printf("csv export file created at %s", csvPath)
		} else {
			file, err := os.Create(csvPath)
			if err != nil {
				return fmt.Errorf("error creating file: %v", err)
			}
			file.Close()
		}
	}

	// open the export file
	file, err := os.OpenFile(csvPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return fmt.Errorf("error opening the csv export file %s: %v", csvPath, err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// write rows
	invitees := invitee.Invitee{}

	records, err := invitees.Scan(rows)
	if err != nil {
		return fmt.Errorf("failed scanning db rows: %v", err)
	}

	for _, record := range records {
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("error in writing to csv: %v", err)
		}
	}

	// create backup for the new emport
	_ = backupImport(csvPath, "export")

	return nil
}
