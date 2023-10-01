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
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/libsql/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

// The release structure is used to represent a GitHub release.
type Release struct {
	// The release name.
	Name string `json:"name"`

	// The release tag name.
	TagName string `json:"tag_name"`

	// The assets associated with the release.
	Assets []Asset `json:"assets"`

	// Indicates whether the release is a draft.
	Draft bool `json:"draft"`

	// Indicates whether the release is a prerelease.
	PreRelease bool `json:"prerelease"`
}

// Asset represents a downloadable asset associated with a release.
type Asset struct {
	// The URL to download the asset.
	URL string `json:"browser_download_url"`

	// The name of the asset.
	Name string `json:"name"`
}

// Bank represents a bank.
type BankCode struct {
	// Code of the bank
	Code string `json:"code"`

	// Ifsc code of the bank
	Ifsc string `json:"ifsc"`

	// Micr code of the bank
	Micr string `json:"micr"`

	// Iin is issuer identification number
	Iin string `json:"iin"`

	// Type of the bank
	Type string `json:"type"`

	// AchCredit is automated clearing house credit is an electronic payment from one bank account to another.
	// ACH credit transfers include direct deposit, payroll and vendor payments.
	// AchCredit specifies whether the bank supports ACH credit
	AchCredit bool `json:"ach_credit"`

	// AchDebit is automated clearing house debit is an electronic payment from one bank account to another.
	// ACH debit transfers include consumer payments on insurance premiums, mortgage loans, and other kinds of bills.
	// AchDebit specifies whether the bank supports ACH debit
	AchDebit bool `json:"ach_debit"`

	// Apbs is Aadhaar Payments Bridge System is a system developed by NPCI for DBT payments.
	// Apbs specifies whether the bank supports APBS
	Apbs bool `json:"apbs"`

	// NachCredit is National Automated Clearing House is a centralized system launched with an aim to
	// consolidate multiple ECS systems running across the country and provides a framework
	// for the harmonization of standard & practices.
	// NachCredit specifies whether the bank supports NACH credit
	NachDebit bool `json:"nach_debit"`
}

//nolint:funlen,gocyclo
func main() {
	ctx := context.Background()
	url := "https://api.github.com/repos/razorpay/ifsc/releases/latest"

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
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

	var release Release
	err = json.NewDecoder(resp.Body).Decode(&release)
	if err != nil {
		slog.ErrorContext(ctx, "error decoding response", slog.Any("err", err), slog.String("url", url))
		return
	}

	if len(release.Assets) == 0 {
		panic("No assets found")
	}

	const perm = 0644
	// Read the TagName from version.txt file
	f, err := os.OpenFile("./tools/finly/version.txt", os.O_RDWR|os.O_CREATE, perm)
	if err != nil {
		slog.ErrorContext(ctx, "error opening file", slog.Any("err", err), slog.String("file", "./tools/finly/version.txt"))
		return
	}

	var TagName string
	_, err = fmt.Fscanf(f, "TAG_VERSION=%s", &TagName)
	if err != nil {
		slog.ErrorContext(ctx, "error scanning file", slog.Any("err", err), slog.String("file", "./tools/finly/version.txt"))
		return
	}

	if release.TagName == TagName || release.Draft || release.PreRelease {
		slog.InfoContext(ctx, "no update required", slog.String("current_version", TagName), slog.String("latest_version", release.TagName))
		return
	}

	var assetURL, assetName, bankURL string
	for _, asset := range release.Assets {
		if asset.Name == "IFSC.csv" {
			assetURL = asset.URL
			assetName = asset.Name
			break
		}

		if asset.Name == "banks.json" {
			bankURL = asset.URL
		}
	}

	if assetURL == "" {
		slog.ErrorContext(ctx, "no asset found", slog.String("asset_name", "IFSC.csv"))
		return
	}

	req, err = http.NewRequestWithContext(context.TODO(), "GET", assetURL, http.NoBody)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", slog.Any("err", err), slog.String("url", assetURL))
		return
	}

	resp, err = client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "error making request", slog.Any("err", err), slog.String("url", assetURL))
		return
	}
	defer resp.Body.Close()

	// Open the SQLite database
	db, err := sql.Open("libsql", "file:./finly/data/finly.db")
	if err != nil {
		slog.ErrorContext(ctx, "error opening database", slog.Any("err", err))
		return
	}
	defer db.Close()

	// Parse the CSV data from the HTTP request body
	reader := csv.NewReader(resp.Body)
	reader.FieldsPerRecord = -1 // Allow variable number of fields per record
	reader.TrimLeadingSpace = true

	// Begin a transaction
	tx, err := db.Begin()
	if err != nil {
		slog.ErrorContext(ctx, "error beginning transaction", slog.Any("err", err))
		return
	}
	defer func() {
		if err != nil {
			err = tx.Rollback()
			if err != nil {
				slog.ErrorContext(ctx, "error rolling back transaction", slog.Any("err", err))
				return
			}
		}
	}()

	_, err = tx.Exec("DROP TABLE IF EXISTS bank")
	if err != nil {
		slog.ErrorContext(ctx, "error dropping table", slog.Any("err", err))
		return
	}

	_, err = tx.Exec("DROP TABLE IF EXISTS version")
	if err != nil {
		slog.ErrorContext(ctx, "error dropping table", slog.Any("err", err))
		return
	}

	// Create the Bank table if it doesn't exist
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS bank (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		code TEXT,
		ifsc TEXT UNIQUE,
		branch TEXT,
		center TEXT,
		district TEXT,
		state TEXT,
		address TEXT,
		contact TEXT,
		imps BOOLEAN,
		rtgs BOOLEAN,
		city TEXT,
		iso3166 TEXT,
		neft BOOLEAN,
		micr TEXT,
		upi BOOLEAN,
		swift TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    )`)
	if err != nil {
		slog.ErrorContext(ctx, "error creating table", slog.Any("err", err))
		return
	}

	// Create the version table if it doesn't exist
	_, err = tx.Exec(`CREATE TABLE IF NOT EXISTS version (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		version TEXT UNIQUE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		slog.ErrorContext(ctx, "error creating table", slog.Any("err", err))
		return
	}

	// Prepare the insert statement
	stmt, err := tx.Prepare(`INSERT INTO bank (
		name,
		ifsc,
		branch,
		center,
		district,
		state,
		address,
		contact,
		imps,
		rtgs,
		city,
		iso3166,
		neft,
		micr,
		upi,
		swift,
		code
    ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		slog.ErrorContext(ctx, "error preparing statement", slog.Any("err", err))
		return
	}
	defer stmt.Close()

	// Insert the parsed data into the Bank table
	var values []interface{}
	_, err = reader.Read()
	if err != nil {
		slog.ErrorContext(ctx, "error reading csv", slog.Any("err", err))
		return
	}

	req, err = http.NewRequestWithContext(context.TODO(), "GET", bankURL, http.NoBody)
	if err != nil {
		slog.ErrorContext(ctx, "error creating request", slog.Any("err", err), slog.String("url", bankURL))
		return
	}

	resp, err = client.Do(req)
	if err != nil {
		slog.ErrorContext(ctx, "error making request", slog.Any("err", err), slog.String("url", bankURL))
		return
	}

	defer resp.Body.Close()

	var banks map[string]BankCode
	err = json.NewDecoder(resp.Body).Decode(&banks)
	if err != nil {
		slog.ErrorContext(ctx, "error decoding response", slog.Any("err", err), slog.Any("response_body", resp.Body))
		return
	}

	ifscBankCode := make(map[string]string)
	for _, bank := range banks {
		ifscBankCode[bank.Ifsc] = bank.Code
	}

	var record []string
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			slog.ErrorContext(ctx, "error reading csv", slog.Any("err", err))
			return
		}

		values = make([]interface{}, len(record))

		for i, v := range record {
			values[i] = v
		}

		bankCode, ok := ifscBankCode[record[1]]
		if !ok {
			values = append(values, record[1][:4])
		} else {
			values = append(values, bankCode)
		}
		_, err = stmt.Exec(values...)
		if err != nil {
			slog.ErrorContext(ctx, "error executing statement", slog.Any("err", err))
			return
		}
	}

	stmt, err = tx.Prepare(`INSERT INTO version (
		name,
		version
	) VALUES (?, ?)`)
	if err != nil {
		slog.ErrorContext(ctx, "error preparing statement", slog.Any("err", err))
		return
	}

	_, err = stmt.Exec(assetName, release.TagName)
	if err != nil {
		slog.ErrorContext(ctx, "error executing statement", slog.Any("err", err))
		return
	}

	// Commit the transaction
	err = tx.Commit()
	if err != nil {
		slog.ErrorContext(ctx, "error committing transaction", slog.Any("err", err))
		return
	}

	// Update the TAG_VERSION in version.txt file
	_, err = f.Seek(0, 0)
	if err != nil {
		slog.ErrorContext(ctx, "error seeking file", slog.Any("err", err), slog.String("file", "./tool/finly/version.txt"))
		return
	}
	_, err = fmt.Fprintf(f, "TAG_VERSION=%s", release.TagName)
	if err != nil {
		slog.ErrorContext(ctx, "error writing to file", slog.Any("err", err), slog.String("file", "./tool/finly/version.txt"))
		return
	}
	err = f.Close()
	if err != nil {
		slog.ErrorContext(ctx, "error closing file", slog.Any("err", err), slog.String("file", "./tool/finly/version.txt"))
		return
	}

	slog.InfoContext(ctx, "update successful", slog.String("current_version", TagName), slog.String("latest_version", release.TagName))
}
