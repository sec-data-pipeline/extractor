package api

import "errors"

func GetNewFilings(cik string, filingIDs []string) ([]*Filing, error) {
	filingsData, err := getFilingsData(cik)
	if err != nil {
		return nil, err
	}
	filings := transformFilings(filingsData)
	return missingFilings(filings, filingIDs), nil
}

func GetMainFile(cik string, flng *Filing) (*File, error) {
	filesData, err := getFilesData(cik, flng.SECID)
	if err != nil {
		return nil, err
	}
	indexFile := flng.rawID + "-index.html"
	content, err := getFileContent(cik, flng.SECID, indexFile)
	if err != nil {
		return nil, err
	}
	mainFile, err := getMainFileName(content)
	if err != nil {
		return nil, err
	}
	files := transformFiles(filesData)
	for _, f := range files {
		if f.Name != mainFile {
			continue
		}
		f.Content, err = getFileContent(cik, flng.SECID, f.Name)
		if err != nil {
			return nil, err
		}
		f.Extension, err = getFileExtension(f.Name)
		return f, nil
	}
	return nil, errors.New("Main file not found")
}
