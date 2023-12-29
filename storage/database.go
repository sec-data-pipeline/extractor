package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sec-data-pipeline/extractor/models"
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

func (db *Database) GetCompanies() ([]models.Company, error) {
	stmt := `SELECT id, cik FROM company;`
	rows, err := db.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var companies []models.Company
	var tmp models.Company
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

func (db *Database) GetLatestFiling() (*models.Filing, error) {
	stmt := `SELECT company.cik, filing.id, filing.sec_id 
	FROM company, filing 
	WHERE company.id = filing.company_id 
	ORDER BY id DESC LIMIT 1;`
	row := db.QueryRow(stmt)
	var company models.Company
	filing := &models.Filing{Company: &company}
	if err := row.Scan(&company.CIK, &filing.ID, &filing.SECID); err != nil {
		return nil, err
	}
	return filing, nil
}

func (db *Database) GetFilings(company *models.Company) ([]models.Filing, error) {
	stmt := `SELECT filing.sec_id FROM filing, company 
	WHERE company.id = filing.company_id AND company.cik = $1;`
	rows, err := db.Query(stmt, company.CIK)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var filings []models.Filing
	var tmp models.Filing
	for rows.Next() {
		if err := rows.Scan(&tmp.SECID); err != nil {
			return nil, err
		}
		filings = append(filings, tmp)
	}
	return filings, nil
}

func (db *Database) InsertFiling(filing *models.Filing) error {
	stmt := `INSERT INTO filing (
		company_id,
		sec_id,
		form,
		size,
		filing_date,
		report_date,
		acceptance_date	
	) VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING id;`
	row := db.QueryRow(
		stmt,
		filing.Company.ID,
		filing.SECID,
		filing.Form,
		filing.Size,
		filing.FilingDate,
		filing.ReportDate,
		filing.AcceptanceDate,
	)
	var filingID int = -1
	if err := row.Scan(&filingID); err != nil {
		return err
	}
	filing.ID = filingID
	return nil
}

func (db *Database) GetFiles(filing *models.Filing) ([]models.File, error) {
	stmt := `SELECT file.name FROM file, filing 
	WHERE filing.ID = file.filing_id AND filing.ID = $1;`
	rows, err := db.Query(stmt, filing.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var files []models.File
	var tmp models.File
	for rows.Next() {
		if err := rows.Scan(&tmp.Name); err != nil {
			return nil, err
		}
		files = append(files, tmp)
	}
	return files, nil
}

func (db *Database) InsertFile(file *models.File) error {
	stmt := `INSERT INTO file (filing_id, name, size, last_modified) VALUES ($1, $2, $3, $4);`
	_, err := db.Exec(stmt, file.Filing.ID, file.Name, file.Size, file.LastModified)
	if err != nil {
		return err
	}
	return nil
}
