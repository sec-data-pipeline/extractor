package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sec-data-pipeline/extractor/filing"
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

func (db *Database) GetCompanies() ([]filing.Company, error) {
	stmt := `SELECT id, cik FROM company;`
	rows, err := db.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var companies []filing.Company
	var tmp filing.Company
	for rows.Next() {
		if err := rows.Scan(&tmp.ID, &tmp.CIK); err != nil {
			return nil, err
		}
		companies = append(companies, tmp)
	}
	return companies, nil
}

func (db *Database) GetFilingCount() (int, error) {
	stmt := `SELECT COUNT(*) FROM filing;`
	row := db.QueryRow(stmt)
	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (db *Database) GetLatestFiling() (*filing.Filing, error) {
	stmt := `SELECT company.cik, filing.id, filing.sec_id 
	FROM company, filing 
	WHERE company.id = filing.company_id 
	ORDER BY id DESC LIMIT 1;`
	row := db.QueryRow(stmt)
	var company filing.Company
	filing := &filing.Filing{Company: &company}
	if err := row.Scan(&company.CIK, &filing.ID, &filing.SECID); err != nil {
		return nil, err
	}
	return filing, nil
}

func (db *Database) GetFilings(company *filing.Company) ([]filing.Filing, error) {
	stmt := `SELECT filing.sec_id FROM filing, company 
	WHERE company.id = filing.company_id AND company.cik = $1;`
	rows, err := db.Query(stmt, company.CIK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var filings []filing.Filing
	var tmp filing.Filing
	for rows.Next() {
		if err := rows.Scan(&tmp.SECID); err != nil {
			return nil, err
		}
		filings = append(filings, tmp)
	}
	return filings, nil
}

func (db *Database) InsertFiling(filing *filing.Filing) error {
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
		filing.Company.ID,
		filing.SECID,
		filing.Form,
		filing.File.Name,
		filing.FilingDate,
		filing.ReportDate,
		filing.AcceptanceDate,
		filing.File.LastModified,
	)
	if err != nil {
		return err
	}
	return nil
}
