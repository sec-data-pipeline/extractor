package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type Database struct {
	*sql.DB
}

type DBConnParams struct {
	DBHost string
	DBPort string
	DBName string
	DBUser string
	DBPass string
	SSL    string
}

func NewDB(connParams *DBConnParams) (*Database, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		connParams.DBHost,
		connParams.DBPort,
		connParams.DBUser,
		connParams.DBName,
		connParams.DBPass,
		connParams.SSL,
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Database{db}, nil
}

func (db *Database) GetCompanies() ([]*company, error) {
	stmt := `SELECT id, cik FROM company;`
	rows, err := db.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var companies []*company
	for rows.Next() {
		var tmp company
		if err := rows.Scan(&tmp.ID, &tmp.CIK); err != nil {
			return nil, err
		}
		companies = append(companies, &tmp)
	}
	return companies, nil
}

func (db *Database) GetFilingIDs(cmp *company) ([]string, error) {
	stmt := `SELECT sec_id FROM filing, company 
	WHERE filing.company_id = company.id AND company.id = $1;`
	rows, err := db.Query(stmt, cmp.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ids []string
	for rows.Next() {
		var tmp string
		if err := rows.Scan(&tmp); err != nil {
			return nil, err
		}
		ids = append(ids, tmp)
	}
	return ids, nil
}

func (db *Database) InsertFiling(
	companyID int,
	secID string,
	form string,
	originalFile string,
	filingDate time.Time,
	reportDate time.Time,
	acceptanceDate time.Time,
	lastModified time.Time,
) error {
	stmt := `INSERT INTO filing (
		company_id,
		sec_id,
		form,
		original_file,
		filing_date,
		report_date,
		acceptance_date,
		last_modified_date
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`
	_, err := db.Exec(
		stmt,
		companyID,
		secID,
		form,
		originalFile,
		filingDate,
		reportDate,
		acceptanceDate,
		lastModified,
	)
	if err != nil {
		return err
	}
	return nil
}
