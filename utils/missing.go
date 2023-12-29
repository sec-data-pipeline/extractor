package utils

import (
	"github.com/sec-data-pipeline/extractor/models"
)

func MissingFilings(localFilings []models.Filing, newFilings []models.Filing) []models.Filing {
	missing := []models.Filing{}
outer:
	for _, newFiling := range newFilings {
		for _, localFiling := range localFilings {
			if newFiling.SECID == localFiling.SECID {
				continue outer
			}
		}
		missing = append(missing, newFiling)
	}
	return missing
}

func MissingFiles(localFiles []models.File, newFiles []models.File) []models.File {
	missing := []models.File{}
outer:
	for _, newFile := range newFiles {
		for _, localFile := range localFiles {
			if newFile.Name == localFile.Name {
				continue outer
			}
		}
		missing = append(missing, newFile)
	}
	return missing
}
