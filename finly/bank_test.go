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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetBankByIFSC(t *testing.T) {
	ctx := context.Background()
	testCases := []struct {
		expectedError  error
		expectedOutput *Bank
		name           string
		ifsc           string
	}{
		{
			name:          "valid ifsc",
			ifsc:          "ABHY0065001",
			expectedError: error(nil),
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
			name:          "invalid ifsc",
			ifsc:          "ABHY0069999",
			expectedError: sql.ErrNoRows,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			finly, err := New()
			if err != nil {
				require.NoError(t, err)
			}

			bank, err := finly.GetBankByIFSC(ctx, tc.ifsc)
			if tc.expectedError != nil {
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.EqualValues(t, tc.expectedOutput, bank)
			}
		})
	}
}
