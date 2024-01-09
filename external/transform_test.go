package external

import (
	"database/sql"
	"fmt"
	"testing"
	"time"
)

func TestTransformFilings(t *testing.T) {
	var tests = []struct {
		name  string
		input *filingsResponse
		want  []*Filing
	}{
		{
			"Empty response",
			&filingsResponse{
				Filings: filings{
					Recent: recent{
						AccessNumber: []string{},
						AcceptDate:   []string{},
						FilingDate:   []string{},
						ReportDate:   []string{},
						Form:         []string{},
					},
				},
			},
			[]*Filing{},
		},
		{
			"Only uninteresting filings",
			&filingsResponse{
				Filings: filings{
					Recent: recent{
						AccessNumber: []string{"", "", ""},
						AcceptDate:   []string{"", "", ""},
						FilingDate:   []string{"", "", ""},
						ReportDate:   []string{"", "", ""},
						Form:         []string{"11-K", "10-K/A", "8-Q"},
					},
				},
			},
			[]*Filing{},
		},
		{
			"One interesting filing with valid timestamps",
			&filingsResponse{
				Filings: filings{
					Recent: recent{
						AccessNumber: []string{"FirstFiling-1043984", "", ""},
						AcceptDate:   []string{"2023-12-14T18:15:12.000Z", "", ""},
						ReportDate:   []string{"2023-12-31", "", ""},
						FilingDate:   []string{"2008-01-02", "", ""},
						Form:         []string{"10-K", "10-K/A", "8-Q"},
					},
				},
			},
			[]*Filing{
				{
					secID: "FirstFiling-1043984",
					AcceptDate: sql.NullTime{
						Time:  time.Date(2023, time.December, 14, 18, 15, 12, 0, time.UTC),
						Valid: true,
					},
					ReportDate: sql.NullTime{
						Time:  time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					FilingDate: sql.NullTime{
						Time:  time.Date(2008, time.January, 2, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Form: "10-K",
				},
			},
		},
		{
			"Multiple interesting filings among uninteresting filings",
			&filingsResponse{
				Filings: filings{
					Recent: recent{
						AccessNumber: []string{
							"FirstFiling-1043984",
							"",
							"",
							"FourthFiling-1043984",
							"FifthFiling-1043984",
							"SixthFiling-1043984",
						},
						AcceptDate: []string{
							"2023-12-14T18:15:12.000Z",
							"",
							"",
							"2023-12-14T18:15:12.000Z",
							"2023-12-14T18:15:12.000Z",
							"2023-12-14T18:15:12.000Z",
						},
						ReportDate: []string{
							"2023-12-31",
							"",
							"",
							"2023-12-31",
							"2023-12-31",
							"2023-12-31",
						},
						FilingDate: []string{
							"2008-01-02",
							"",
							"",
							"2008-01-02",
							"2008-01-02",
							"2008-01-02",
						},
						Form: []string{"10-K", "10-K/A", "8-Q", "10-Q", "8-Q", "10-K"},
					},
				},
			},
			[]*Filing{
				{
					secID: "FirstFiling-1043984",
					AcceptDate: sql.NullTime{
						Time:  time.Date(2023, time.December, 14, 18, 15, 12, 0, time.UTC),
						Valid: true,
					},
					ReportDate: sql.NullTime{
						Time:  time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					FilingDate: sql.NullTime{
						Time:  time.Date(2008, time.January, 2, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Form: "10-K",
				},
				{
					secID: "FourthFiling-1043984",
					AcceptDate: sql.NullTime{
						Time:  time.Date(2023, time.December, 14, 18, 15, 12, 0, time.UTC),
						Valid: true,
					},
					ReportDate: sql.NullTime{
						Time:  time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					FilingDate: sql.NullTime{
						Time:  time.Date(2008, time.January, 2, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Form: "10-Q",
				},
				{
					secID: "SixthFiling-1043984",
					AcceptDate: sql.NullTime{
						Time:  time.Date(2023, time.December, 14, 18, 15, 12, 0, time.UTC),
						Valid: true,
					},
					ReportDate: sql.NullTime{
						Time:  time.Date(2023, time.December, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					FilingDate: sql.NullTime{
						Time:  time.Date(2008, time.January, 2, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					Form: "10-K",
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filings := transformFilings(test.input)
			for i, got := range filings {
				if got.secID != test.want[i].secID {
					t.Errorf("got: %s, want: %s", got.secID, test.want[i].secID)
				}
				msg, ok := checkNullTime(got.AcceptDate, test.want[i].AcceptDate)
				if !ok {
					t.Errorf(fmt.Sprintf("for filing: %s, AcceptDate ", got.secID) + msg)
				}
				msg, ok = checkNullTime(got.FilingDate, test.want[i].FilingDate)
				if !ok {
					t.Errorf(fmt.Sprintf("for filing: %s, FilingDate ", got.secID) + msg)
				}
				msg, ok = checkNullTime(got.ReportDate, test.want[i].ReportDate)
				if !ok {
					t.Errorf(fmt.Sprintf("for filing: %s, ReportDate ", got.secID) + msg)
				}
				if got.Form != test.want[i].Form {
					t.Errorf(
						"for filing: %s, Form got: %s, want: %s",
						got.secID,
						got.Form,
						test.want[i].Form,
					)
				}
			}
		})
	}
}

func TestTransformFiles(t *testing.T) {
	var tests = []struct {
		name  string
		input *filesResponse
		want  []*file
	}{
		{
			"One file with empty timestamp",
			&filesResponse{Dir: directory{Items: []item{{Name: "Test", LastModified: ""}}}},
			[]*file{{Name: "Test", LastModified: sql.NullTime{Valid: false}}},
		},
		{
			"One file with valid timestamp",
			&filesResponse{
				Dir: directory{
					Items: []item{{Name: "Test", LastModified: "2008-03-06 04:20:59"}},
				},
			},
			[]*file{
				{
					Name: "Test",
					LastModified: sql.NullTime{
						Time:  time.Date(2008, time.March, 6, 4, 20, 59, 0, time.UTC),
						Valid: true,
					},
				},
			},
		},
		{
			"One file with invalid timestamp",
			&filesResponse{
				Dir: directory{
					Items: []item{{Name: "Test", LastModified: "2023-12-14T18:15:12.000Z"}},
				},
			},
			[]*file{{Name: "Test", LastModified: sql.NullTime{Valid: false}}},
		},
		{
			"Multiple files with empty timestamps",
			&filesResponse{
				Dir: directory{
					Items: []item{
						{Name: "Test", LastModified: ""},
						{Name: "SecondFile", LastModified: ""},
						{Name: "cv.pdf", LastModified: ""},
					},
				},
			},
			[]*file{
				{Name: "Test", LastModified: sql.NullTime{Valid: false}},
				{Name: "SecondFile", LastModified: sql.NullTime{Valid: false}},
				{Name: "cv.pdf", LastModified: sql.NullTime{Valid: false}},
			},
		},
		{
			"Multiple files with valid timestamps",
			&filesResponse{
				Dir: directory{
					Items: []item{
						{Name: "Test", LastModified: "2008-03-06 04:20:59"},
						{Name: "SecondFile", LastModified: "2025-05-12 12:20:09"},
						{Name: "cv.pdf", LastModified: "2028-07-30 05:10:00"},
					},
				},
			},
			[]*file{
				{
					Name: "Test",
					LastModified: sql.NullTime{
						Time:  time.Date(2008, time.March, 6, 4, 20, 59, 0, time.UTC),
						Valid: true,
					},
				},
				{
					Name: "SecondFile",
					LastModified: sql.NullTime{
						Time:  time.Date(2025, time.May, 12, 12, 20, 9, 0, time.UTC),
						Valid: true,
					},
				},
				{
					Name: "cv.pdf",
					LastModified: sql.NullTime{
						Time:  time.Date(2028, time.July, 30, 5, 10, 0, 0, time.UTC),
						Valid: true,
					},
				},
			},
		},
		{
			"Multiple files with invalid timestamps",
			&filesResponse{
				Dir: directory{
					Items: []item{
						{Name: "Test", LastModified: "2023-12-14T18:15:12.000Z"},
						{Name: "SecondFile", LastModified: "2006-01-02"},
						{Name: "cv.pdf", LastModified: "2026-11-30"},
					},
				},
			},
			[]*file{
				{Name: "Test", LastModified: sql.NullTime{Valid: false}},
				{Name: "SecondFile", LastModified: sql.NullTime{Valid: false}},
				{Name: "cv.pdf", LastModified: sql.NullTime{Valid: false}},
			},
		},
		{
			"Multiple files with valid and invalid timestamps",
			&filesResponse{
				Dir: directory{
					Items: []item{
						{Name: "Test", LastModified: "2023-12-14T18:15:12.000Z"},
						{Name: "SecondFile", LastModified: "2025-05-12 12:20:09"},
						{Name: "cv.pdf", LastModified: "2028-07-30 05:10:00"},
						{Name: "FourthFile", LastModified: ""},
					},
				},
			},
			[]*file{
				{Name: "Test", LastModified: sql.NullTime{Valid: false}},
				{
					Name: "SecondFile",
					LastModified: sql.NullTime{
						Time:  time.Date(2025, time.May, 12, 12, 20, 9, 0, time.UTC),
						Valid: true,
					},
				},
				{
					Name: "cv.pdf",
					LastModified: sql.NullTime{
						Time:  time.Date(2028, time.July, 30, 5, 10, 0, 0, time.UTC),
						Valid: true,
					},
				},
				{Name: "FourthFile", LastModified: sql.NullTime{Valid: false}},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			files := transformFiles(test.input)
			for i, got := range files {
				if got.Name != test.want[i].Name {
					t.Errorf("got: %s, want: %s", got.Name, test.want[i].Name)
				}
				msg, ok := checkNullTime(got.LastModified, test.want[i].LastModified)
				if !ok {
					t.Errorf(fmt.Sprintf("for file: %s, ", got.Name) + msg)
				}
			}
		})
	}
}

func TestParseNullTime(t *testing.T) {
	var tests = []struct {
		name   string
		layout string
		value  string
		want   sql.NullTime
	}{
		{"Full layout empty string", "2006-01-02 15:04:05", "", sql.NullTime{Valid: false}},
		{"Short layout empty string", "2006-01-02", "", sql.NullTime{Valid: false}},
		{"RFC3339 layout empty string", time.RFC3339, "", sql.NullTime{Valid: false}},
		{
			"Full layout valid string",
			"2006-01-02 15:04:05",
			"2004-09-10 16:47:30",
			sql.NullTime{
				Time:  time.Date(2004, time.September, 10, 16, 47, 30, 0, time.UTC),
				Valid: true,
			},
		},
		{
			"Short layout valid string",
			"2006-01-02",
			"2023-08-15",
			sql.NullTime{Time: time.Date(2023, time.August, 15, 0, 0, 0, 0, time.UTC), Valid: true},
		},
		{
			"RFC3339 layout valid string",
			time.RFC3339,
			"2023-12-14T18:15:12.000Z",
			sql.NullTime{
				Time:  time.Date(2023, time.December, 14, 18, 15, 12, 0, time.UTC),
				Valid: true,
			},
		},
		{
			"Full layout invalid string",
			"2006-01-02 15:04:05",
			"2023-08-15",
			sql.NullTime{
				Valid: false,
			},
		},
		{
			"Short layout invalid string",
			"2006-01-02",
			"2004-09-10 16:47:30",
			sql.NullTime{Valid: false},
		},
		{"RFC3339 layout invalid string", time.RFC3339, "2023-08-15", sql.NullTime{Valid: false}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := parseNullTime(test.layout, test.value)
			msg, ok := checkNullTime(got, test.want)
			if !ok {
				t.Errorf(msg)
			}
		})
	}
}

func checkNullTime(got sql.NullTime, want sql.NullTime) (string, bool) {
	if got.Valid != want.Valid {
		return fmt.Sprintf("got validity: %t, want: %t", got.Valid, want.Valid), false
	}
	if got.Time != want.Time {
		return fmt.Sprintf("got: %s, want: %s", got.Time.String(), want.Time.String()), false
	}
	return "", true
}
