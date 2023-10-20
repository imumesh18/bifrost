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

// atlas package provides a store that interacts with the database to retrieve geo location information.
package atlas

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	_ "github.com/libsql/libsql-client-go/libsql" // Import the libsql driver
	_ "modernc.org/sqlite"                        // Import the sqlite driver
)

const getGeoLocationByPostalCodeQuery = `SELECT country_code, postal_code, place_name,
admin_name1, admin_code1, admin_name2,
admin_code2, admin_name3, admin_code3,
latitude, longitude, accuracy
FROM geo_location WHERE postal_code = ?`

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

var ErrGeoLocationNotFound = errors.New("geo location not found")

// Atlas represents a geographic location service.
type Atlas struct {
	db *sql.DB
}

// New creates a new instance of Atlas and initializes the database connection.
// It returns a pointer to Atlas and an error if the connection fails.
func New() (*Atlas, error) {
	db, err := sql.Open("libsql", "file:"+getDBPath()+"?mode=ro")
	if err != nil {
		return nil, err
	}

	return &Atlas{db: db}, nil
}

func getDBPath() string {
	// Get the absolute path to the current working directory
	wd, err := findProjectRoot()
	if err != nil {
		panic(err)
	}

	// Join the atlas directory path with the database file name
	return filepath.Join(wd, "atlas", "data", "atlas.db")
}

func findProjectRoot() (string, error) {
	// Start with the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse upwards to find a marker file or directory that indicates the project root
	for {
		// Check if a specific file or directory exists in the current directory
		// that indicates the project root (e.g., a marker file like 'go.mod')
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}

		// If not found, move up one directory level
		parentDir := filepath.Dir(currentDir)
		// If we're already at the root directory, stop searching
		if parentDir == currentDir {
			return "", nil
		}
		currentDir = parentDir
	}
}

// GetGeoLocationByPostalCode retrieves a GeoLocation struct from the database by postal code.
// If the GeoLocation is not found, returns ErrGeoLocationNotFound.
func (a *Atlas) GetGeoLocationByPostalCode(ctx context.Context, postalCode string) (*GeoLocation, error) {
	var geo GeoLocation

	err := a.db.QueryRowContext(ctx, getGeoLocationByPostalCodeQuery, postalCode).Scan(
		&geo.CountryCode,
		&geo.PostalCode,
		&geo.PlaceName,
		&geo.AdminName1,
		&geo.AdminCode1,
		&geo.AdminName2,
		&geo.AdminCode2,
		&geo.AdminName3,
		&geo.AdminCode3,
		&geo.Latitude,
		&geo.Longitude,
		&geo.Accuracy,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrGeoLocationNotFound
		}
		return nil, err
	}

	return &geo, nil
}
