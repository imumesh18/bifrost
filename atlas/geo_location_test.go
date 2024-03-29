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

package atlas

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetGeoLocationByPostalCode(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		expectedError  error
		expectedOutput *GeoLocation
		mockDB         func(mock sqlmock.Sqlmock)
		name           string
		postalCode     string
	}{
		{
			name:       "valid postal code",
			postalCode: "400001",
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(getGeoLocationByPostalCodeQuery).WithArgs("400001").WillReturnRows(sqlmock.NewRows([]string{
					"country_code",
					"postal_code",
					"place_name",
					"admin_name_1",
					"admin_code_1",
					"admin_name_2",
					"admin_code_2",
					"admin_name_3",
					"admin_code_3",
					"latitude",
					"longitude",
					"accuracy",
				}).AddRow(
					"IN",
					"560095",
					"Koramangala VI Bk",
					"Karnataka",
					"19",
					"Bengaluru",
					"583",
					"Bangalore South",
					"",
					13.1077,
					77.581,
					1,
				))
			},
			expectedError: error(nil),
			expectedOutput: &GeoLocation{
				CountryCode: "IN",
				PostalCode:  "560095",
				PlaceName:   "Koramangala VI Bk",
				AdminName1:  "Karnataka",
				AdminCode1:  "19",
				AdminName2:  "Bengaluru",
				AdminCode2:  "583",
				AdminName3:  "Bangalore South",
				AdminCode3:  "",
				Latitude:    13.1077,
				Longitude:   77.581,
				Accuracy:    1,
			},
		},
		{
			name:       "invalid postal code",
			postalCode: "999999",
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(getGeoLocationByPostalCodeQuery).WithArgs("999999").WillReturnError(sql.ErrNoRows)
			},
			expectedError: ErrGeoLocationNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)
			atlas := Atlas{
				db: db,
			}
			tc.mockDB(mock)

			geoLocation, err := atlas.GetGeoLocationByPostalCode(ctx, tc.postalCode)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expectedOutput, geoLocation)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
