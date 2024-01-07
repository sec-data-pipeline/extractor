package storage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

type company struct {
	ID  int
	CIK string
}

type Database interface {
	GetCompanies() ([]*company, error)
	GetFilingIDs(cmpID int) ([]string, error)
	InsertFiling(
		cmpID int,
		secID string,
		form string,
		ogFile string,
		filDate time.Time,
		repDate time.Time,
		acptDate time.Time,
		lMDate time.Time,
	) error
}

type postgresDB struct {
	*sql.DB
}

type postgresConnParams struct {
	DBHost string
	DBPort string
	DBName string
	DBUser string
	DBPass string
	SSL    string
}

func NewPostgresConn(connParams *postgresConnParams) (*postgresDB, error) {
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
	return &postgresDB{db}, nil
}

func (db *postgresDB) GetCompanies() ([]*company, error) {
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

func (db *postgresDB) GetFilingIDs(cmpID int) ([]string, error) {
	stmt := `SELECT sec_id FROM filing, company 
	WHERE filing.company_id = company.id AND company.id = $1;`
	rows, err := db.Query(stmt, cmpID)
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

func (db *postgresDB) InsertFiling(
	cmpID int,
	secID string,
	form string,
	ogFile string,
	filDate time.Time,
	repDate time.Time,
	acptDate time.Time,
	lMDate time.Time,
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
	_, err := db.Exec(stmt, cmpID, secID, form, ogFile, filDate, repDate, acptDate, lMDate)
	if err != nil {
		return err
	}
	return nil
}
