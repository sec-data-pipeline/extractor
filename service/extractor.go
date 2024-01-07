package service

import (
	"github.com/sec-data-pipeline/filing-extractor/external"
	"github.com/sec-data-pipeline/filing-extractor/storage"
)

type Extractor struct {
	api     external.API
	db      storage.Database
	archive storage.FileStorage
	logger  storage.Logger
}

func NewExtractorService(
	api external.API,
	db storage.Database,
	archive storage.FileStorage,
	logger storage.Logger,
) *Extractor {
	return &Extractor{api: api, db: db, archive: archive, logger: logger}
}

func (s *Extractor) Run() error {
	companies, err := s.db.GetCompanies()
	if err != nil {
		return err
	}
	for _, cmp := range companies {
		filIDs, err := s.db.GetFilingIDs(cmp.ID)
		if err != nil {
			return err
		}
		filings, err := s.getMissingFilings(cmp.CIK, filIDs)
		if err != nil {
			s.logger.Log(err.Error())
			continue
		}
		for _, fil := range filings {
			mainFile, err := s.api.GetMainFile(cmp.CIK, fil)
			if err != nil {
				s.logger.Log(err.Error())
				continue
			}
			ex, err := mainFile.GetExtension()
			if err != nil {
				s.logger.Log(err.Error())
				continue
			}
			err = s.db.InsertFiling(
				cmp.ID,
				fil.GetID(),
				fil.Form,
				mainFile.Name,
				fil.FilingDate,
				fil.ReportDate,
				fil.AcceptanceDate,
				mainFile.LastModified,
			)
			if err != nil {
				s.logger.Log(err.Error())
				continue
			}
			err = s.archive.PutObject(fil.GetID()+ex, mainFile.Content)
			if err != nil {
				s.logger.Log(err.Error())
			}
		}
	}
	return nil
}

func (s *Extractor) getMissingFilings(cik string, got []string) ([]*external.Filing, error) {
	filings, err := s.api.GetFilings(cik)
	if err != nil {
		return nil, err
	}
	var missing []*external.Filing
outer:
	for _, fil := range filings {
		for _, id := range got {
			if fil.GetID() == id {
				continue outer
			}
		}
		missing = append(missing, fil)
	}
	return missing, nil
}
