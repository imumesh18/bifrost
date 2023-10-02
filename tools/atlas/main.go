// Copyright (C) 2023 Umesh Yadav
//
// Licensed under the MIT License (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://opensource.org/licenses/MIT
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"archive/zip"
	"context"
	"database/sql"
	"encoding/csv"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

// GeoLocation represents a geographical location details
type GeoLocation struct {
	// ISO country code abbreviation
	CountryCode string `json:"country_code,omitempty"`

	// Postal code or zip code
	PostalCode string `json:"postal_code,omitempty"`

	// Name of the place
	PlaceName string `json:"place_name,omitempty"`

	// First-order administrative division (state, province, region, etc.)
	AdminName1 string `json:"admin_name_1,omitempty"`

	// Code for the first-order administrative division
	AdminCode1 string `json:"admin_code_1,omitempty"`

	// Second-order administrative division (county, district, etc.)
	AdminName2 string `json:"admin_name_2,omitempty"`

	// Code for the second-order administrative division
	AdminCode2 string `json:"admin_code_2,omitempty"`

	// Third-order administrative division (township, municipality, etc.)
	AdminName3 string `json:"admin_name_3,omitempty"`

	// Code for the third-order administrative division
	AdminCode3 string `json:"admin_code_3,omitempty"`

	// Latitude of the location
	Latitude float64 `json:"latitude,omitempty"`

	// Longitude of the location
	Longitude float64 `json:"longitude,omitempty"`

	// Accuracy of the location in meters
	Accuracy int `json:"accuracy,omitempty"`
}

//nolint:funlen,gocyclo
func main() {
	ctx := context.Background()
	// Download the allCountries.zip file
	url := "http://download.geonames.org/export/zip/allCountries.zip"
	req, err := http.NewRequestWithContext(context.TODO(), "GET", url, http.NoBody)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", slog.Any("err", err), slog.String("url", url))
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "error making request", slog.Any("err", err), slog.String("url", url))
		return
	}
	defer resp.Body.Close()

	// Create a temporary file to store the downloaded ZIP archive
	tmpFile, err := os.CreateTemp("", "allCountries-*.zip")
	if err != nil {
		slog.ErrorContext(ctx, "error creating temporary file", slog.Any("err", err))
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write the ZIP archive to the temporary file
	_, err = io.Copy(tmpFile, resp.Body)
	if err != nil {
		slog.ErrorContext(ctx, "error writing to temporary file", slog.Any("err", err))
		return
	}

	// Extract the allCountries.txt file from the ZIP archive
	zipReader, err := zip.OpenReader(tmpFile.Name())
	if err != nil {
		slog.ErrorContext(ctx, "error opening ZIP archive", slog.Any("err", err))
		return
	}
	defer zipReader.Close()

	var allCountriesFile *zip.File
	for _, file := range zipReader.File {
		if file.Name == "allCountries.txt" {
			allCountriesFile = file
			break
		}
	}
	if allCountriesFile == nil {
		slog.ErrorContext(ctx, "allCountries.txt not found in ZIP archive")
		return
	}

	allCountriesReader, err := allCountriesFile.Open()
	if err != nil {
		slog.ErrorContext(ctx, "error opening allCountries.txt", slog.Any("err", err))
		return
	}
	defer allCountriesReader.Close()

	// Open the SQLite3 database
	db, err := sql.Open("libsql", "file:./atlas/data/atlas.db")
	if err != nil {
		slog.ErrorContext(ctx, "error opening database", slog.Any("err", err))
		return
	}
	defer db.Close()

	// Create the geonames table if it doesn't exist
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS geo_location (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
        country_code TEXT,
        postal_code TEXT,
        place_name TEXT,
        admin_name1 TEXT,
        admin_code1 TEXT,
        admin_name2 TEXT,
        admin_code2 TEXT,
        admin_name3 TEXT,
        admin_code3 TEXT,
        latitude REAL,
        longitude REAL,
        accuracy INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		slog.ErrorContext(ctx, "error creating table", slog.Any("err", err))
		return
	}

	// Parse the allCountries.txt file as a CSV
	reader := csv.NewReader(allCountriesReader)
	reader.Comma = '\t'
	reader.LazyQuotes = true

	// Loop through the CSV records and insert them into the database
	tx, err := db.Begin()
	if err != nil {
		slog.ErrorContext(ctx, "error beginning transaction", slog.Any("err", err))
		return
	}
	stmt, err := tx.Prepare(`INSERT INTO geo_location (
        country_code,
        postal_code,
        place_name,
        admin_name1,
        admin_code1,
        admin_name2,
        admin_code2,
        admin_name3,
        admin_code3,
        latitude,
        longitude,
        accuracy
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		slog.ErrorContext(ctx, "error preparing statement", slog.Any("err", err))
		return
	}
	defer stmt.Close()

	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			slog.ErrorContext(ctx, "error reading csv", slog.Any("err", err))
			return
		}

		// Parse the record into a GeoName struct
		geoName := GeoLocation{
			CountryCode: record[0],
			PostalCode:  record[1],
			PlaceName:   record[2],
			AdminName1:  record[3],
			AdminCode1:  record[4],
			AdminName2:  record[5],
			AdminCode2:  record[6],
			AdminName3:  record[7],
			AdminCode3:  record[8],
		}

		// Parse the latitude and longitude fields as floats
		geoName.Latitude, err = strconv.ParseFloat(record[9], 64)
		if err != nil {
			slog.ErrorContext(ctx, "error parsing latitude", slog.Any("err", err))
			return
		}
		geoName.Longitude, err = strconv.ParseFloat(record[10], 64)
		if err != nil {
			slog.ErrorContext(ctx, "error parsing longitude", slog.Any("err", err))
			return
		}

		// Parse the accuracy field as an integer
		if strings.TrimSpace(record[11]) != "" {
			geoName.Accuracy, err = strconv.Atoi(strings.TrimSpace(record[11]))
		}
		if err != nil {
			slog.ErrorContext(ctx, "error parsing accuracy", slog.Any("err", err))
			return
		}

		// Insert the GeoName struct into the database
		_, err = stmt.Exec(
			geoName.CountryCode,
			geoName.PostalCode,
			geoName.PlaceName,
			geoName.AdminName1,
			geoName.AdminCode1,
			geoName.AdminName2,
			geoName.AdminCode2,
			geoName.AdminName3,
			geoName.AdminCode3,
			geoName.Latitude,
			geoName.Longitude,
			geoName.Accuracy,
		)
		if err != nil {
			slog.ErrorContext(ctx, "error executing statement", slog.Any("err", err))
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		slog.ErrorContext(ctx, "error committing transaction", slog.Any("err", err))
		return
	}

	slog.InfoContext(ctx, "data generated successfully")
}
