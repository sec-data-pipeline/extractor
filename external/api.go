package external

import (
	"errors"
	"strings"
	"time"
)

type API interface {
	GetFilings(cik string) ([]*Filing, error)
	GetMainFile(cik string, fil *Filing) (*file, error)
}

type Filing struct {
	secID          string
	Form           string
	FilingDate     time.Time
	ReportDate     time.Time
	AcceptanceDate time.Time
}

func (f *Filing) GetID() string {
	return strings.Replace(f.secID, "-", "", -1)
}

type file struct {
	Name         string
	Content      []byte
	LastModified time.Time
}

func (f *file) GetExtension() (string, error) {
	if !strings.Contains(f.Name, ".") {
		return "", errors.New("File extension could not be found")
	}
	result := ""
	for i := len(f.Name) - 1; i >= 0; i-- {
		result = string(f.Name[i]) + result
		if string(f.Name[i]) == "." {
			break
		}
	}
	return result, nil
}

func getFile(files []*file, name string) (*file, error) {
	for _, file := range files {
		if file.Name == name {
			return file, nil
		}
	}
	return nil, errors.New("File not in provided list")
}
