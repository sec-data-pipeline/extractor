package api

import "errors"

func GetNewFilings(cik string, filingIDs []string) ([]*Filing, error) {
	filingsData, err := getFilingsData(cik)
	if err != nil {
		return nil, errors.New(
			"[APIError] CIK: '" + cik + "', could not fetch filings data, " + err.Error(),
		)
	}
	filings := transformFilings(filingsData)
	return missingFilings(filings, filingIDs), nil
}

func GetMainFile(cik string, flng *Filing) (*File, error) {
	filesData, err := getFilesData(cik, flng.SECID)
	if err != nil {
		return nil, errors.New(
			"[APIError] Filing: '" + flng.SECID + "', could not fetch files data, " + err.Error(),
		)
	}
	indexFile := flng.rawID + "-index.html"
	content, err := getFileContent(cik, flng.SECID, indexFile)
	if err != nil {
		return nil, errors.New(
			"[APIError] Filing: '" + flng.SECID + "', content could not be found for index file: '" +
				indexFile + "', " + err.Error(),
		)
	}
	mainFile, err := getMainFileName(content)
	if err != nil {
		return nil, errors.New(
			"[APIError] Filing: '" + flng.SECID + "', name of main file could not be found in: '" +
				indexFile + "', " + err.Error(),
		)
	}
	files := transformFiles(filesData)
	for _, f := range files {
		if f.Name != mainFile {
			continue
		}
		f.Content, err = getFileContent(cik, flng.SECID, f.Name)
		if err != nil {
			return nil, errors.New(
				"[APIError] Filing: '" + flng.SECID + "', content could not be found for main file: '" +
					f.Name + "', " + err.Error(),
			)
		}
		f.Extension, err = getFileExtension(f.Name)
		return f, nil
	}
	return nil, errors.New(
		"[APIError] Filing: '" + flng.SECID + "', main file not found in: '" + indexFile + "'",
	)
}
