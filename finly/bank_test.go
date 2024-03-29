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

package finly

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBankByIFSC(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		expectedError  error
		expectedOutput *Bank
		mockDB         func(mock sqlmock.Sqlmock)
		name           string
		ifsc           string
	}{
		{
			name:          "valid ifsc",
			ifsc:          "ABHY0065001",
			expectedError: error(nil),
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(getBankByIFSCQuery).WithArgs("ABHY0065001").WillReturnRows(sqlmock.NewRows([]string{
					"name",
					"code",
					"ifsc",
					"branch",
					"center",
					"district",
					"state",
					"address",
					"contact",
					"imps",
					"rtgs",
					"city",
					"iso3166",
					"neft",
					"micr",
					"upi",
					"swift",
				}).AddRow(
					"Abhyudaya Co-operative Bank",
					"ABHY",
					"ABHY0065001",
					"Abhyudaya Co-operative Bank IMPS",
					"MUMBAI",
					"MUMBAI",
					"MAHARASHTRA",
					"ABHYUDAYA BUILDING, KAMAL NATH MARG,NEHRU NAGAR,KURLA-EAST,MUMBAI-400024",
					"+919653261383",
					"1",
					"1",
					"MUMBAI",
					"IN-MH",
					"1",
					"400065001",
					"1",
					"",
				))
			},
			expectedOutput: &Bank{
				Name:     "Abhyudaya Co-operative Bank",
				State:    "MAHARASHTRA",
				City:     "MUMBAI",
				Micr:     "400065001",
				Branch:   "Abhyudaya Co-operative Bank IMPS",
				Code:     "ABHY",
				Contact:  "+919653261383",
				Ifsc:     "ABHY0065001",
				District: "MUMBAI",
				Address:  "ABHYUDAYA BUILDING, KAMAL NATH MARG,NEHRU NAGAR,KURLA-EAST,MUMBAI-400024",
				Center:   "MUMBAI",
				Swift:    "",
				Iso3166:  "IN-MH",
				Neft:     true,
				Rtgs:     true,
				Imps:     true,
				Upi:      true,
			},
		},
		{
			name: "invalid ifsc",
			ifsc: "ABHY0069999",
			mockDB: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(getBankByIFSCQuery).WithArgs("ABHY0069999").WillReturnError(sql.ErrNoRows)
			},
			expectedError: ErrBankNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			require.NoError(t, err)

			tc.mockDB(mock)

			finly := &Finly{
				store: db,
			}

			bank, err := finly.GetBankByIFSC(ctx, tc.ifsc)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expectedOutput, bank)
			}

			err = mock.ExpectationsWereMet()
			assert.NoError(t, err)
		})
	}
}
