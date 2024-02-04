package storage

import (
	"database/sql"
	"fmt"

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
		filDate sql.NullTime,
		repDate sql.NullTime,
		acptDate sql.NullTime,
		lMDate sql.NullTime,
	) error
}

type postgresDB struct {
	*sql.DB
}

type postgresParams struct {
	DBHost string `json:"DB_HOST"`
	DBPort string `json:"DB_PORT"`
	DBName string `json:"DB_NAME"`
	DBUser string `json:"DB_USER"`
	DBPass string `json:"DB_PASS"`
	ssl    string
}

func NewPostgres(params *postgresParams) (*postgresDB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		params.DBHost,
		params.DBPort,
		params.DBUser,
		params.DBName,
		params.DBPass,
		params.ssl,
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
	filDate sql.NullTime,
	repDate sql.NullTime,
	acptDate sql.NullTime,
	lMDate sql.NullTime,
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
