package database

import (
	"database/sql"
	"fmt"
	"log"
)

var initDefault = MasterTaxDeduction{
	CycleYear:               2567,
	AmtKReceiptMax:          50000,
	AmtDonationMax:          100000,
	AmtPersonalDeductionMin: 60000,
	CreatedBy:               "SYSTEM_INIT",
	UpdatedBy:               "SYSTEM_INIT",
}

func CreateMasterTaxDeduction(db *sql.DB) error {
	_, err := db.Exec(`
		create table IF NOT EXISTS master_tax_deduction (
			id 							SERIAL 			NOT NULL,
			cycle_year 					SMALLINT 		NOT NULL,
			amt_k_receipt_max 			NUMERIC(12, 2) 	NOT NULL,
			amt_donation_max 			NUMERIC(12, 2) 	NOT NULL,
			amt_personal_deduction_min 	NUMERIC(12, 2) 	NOT NULL,
			created_at 					TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			created_by 					VARCHAR(64) 	NOT NULL,
			updated_at 					TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_by 					VARCHAR(64) 	NOT NULL,
			CONSTRAINT master_tax_deduction_pkey PRIMARY KEY (id),
			CONSTRAINT master_tax_deduction_unq_cycle_year UNIQUE (cycle_year)
		);
	`)
	if err != nil {
		return fmt.Errorf("create table master_tax_deduction: %v", err)
	}

	log.Println("Table master_tax_deduction created!")
	return nil
}

func InsertMasterTaxDeduction(db *sql.DB) error {
	sql := `
		INSERT INTO master_tax_deduction (cycle_year,amt_k_receipt_max,amt_donation_max,amt_personal_deduction_min,created_by,updated_by)
		VALUES ($1,$2,$3,$4,$5,$6)
		ON CONFLICT DO NOTHING;
	`
	// default value
	_, err := db.Exec(
		sql,
		initDefault.CycleYear,
		initDefault.AmtKReceiptMax,
		initDefault.AmtDonationMax,
		initDefault.AmtPersonalDeductionMin,
		initDefault.CreatedBy,
		initDefault.UpdatedBy,
	)

	if err != nil {
		return fmt.Errorf("insert table master_tax_deduction of cycle_year=%v: %v", initDefault.CycleYear, err)
	}

	log.Println("Table master_tax_deduction inserted!")
	return nil
}
