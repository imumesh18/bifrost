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

// Package finly provides a store that interacts with the database to retrieve bank information.
package finly

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/libsql/libsql-client-go/libsql" // Import the libsql driver
	_ "modernc.org/sqlite"                        // Import the sqlite driver
)

// getBankByIFSCQuery is the query to get the bank by ifsc
const getBankByIFSCQuery = `SELECT name, code, ifsc, branch, center,
district, state, address, contact,
imps, rtgs, city, iso3166,
neft, micr, upi, swift
FROM bank WHERE ifsc = ?`

var ErrBankNotFound = errors.New("bank not found")

// Bank entity represents the bank across the indian banking system
type Bank struct {
	// Name specifies the name of the bank
	Name string `json:"name"`
	// State specifics which state the bank is in
	State string `json:"state"`
	// City specifies which city the bank is in
	City string `json:"city"`
	// Micr specifies the micr code of the bank
	Micr string `json:"micr"`
	// Branch specifies the branch of the bank
	Branch string `json:"branch"`
	// Code specifies the code of the bank which is unique and 4 letters long
	Code string `json:"code"`
	// Contact information of the bank
	Contact string `json:"contact"`
	// Ifsc code of the bank branch
	Ifsc string `json:"ifsc"`
	// District specifies which district the bank is in
	District string `json:"district"`
	// Address specifies the address of the bank
	Address string `json:"address"`
	// Center specifies the center of the bank
	Center string `json:"center"`
	// Swift specifies the swift code of the bank
	Swift string `json:"swift"`
	// Iso3166 specifies the iso3166 code of the bank
	Iso3166 string `json:"iso3166"`
	// Neft specifies whether the bank supports neft
	Neft bool `json:"neft"`
	// Rtgs specifies whether the bank supports rtgs
	Rtgs bool `json:"rtgs"`
	// Imps specifies whether the bank supports imps
	Imps bool `json:"imps"`
	// Upi specifies whether the bank supports upi
	Upi bool `json:"upi"`
}

// Finly is the store that interacts with the database
type Finly struct {
	store *sql.DB
}

// New returns a new finly instance that interacts with the database.
// It returns an error if it fails to open the database.
func New() (*Finly, error) {
	db, err := sql.Open("libsql", "file:./data/finly.db")
	if err != nil {
		return nil, err
	}

	return &Finly{store: db}, nil
}

// GetBankByIFSC returns a Bank instance by its IFSC code.
// It returns an error if it fails to query the database.
func (b *Finly) GetBankByIFSC(ctx context.Context, ifsc string) (*Bank, error) {
	var i Bank
	err := b.store.QueryRowContext(ctx, getBankByIFSCQuery, ifsc).Scan(
		&i.Name,
		&i.Code,
		&i.Ifsc,
		&i.Branch,
		&i.Center,
		&i.District,
		&i.State,
		&i.Address,
		&i.Contact,
		&i.Imps,
		&i.Rtgs,
		&i.City,
		&i.Iso3166,
		&i.Neft,
		&i.Micr,
		&i.Upi,
		&i.Swift,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBankNotFound
		}

		return nil, err
	}

	return &i, nil
}
