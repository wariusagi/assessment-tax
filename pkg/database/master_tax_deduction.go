package database

import "time"

type MasterTaxDeduction struct {
	Id                      int       `db:"id"`
	CycleYear               int       `db:"cycle_year"`
	AmtKReceiptMax          float64   `db:"amt_k_receipt_max"`
	AmtDonationMax          float64   `db:"amt_donation_max"`
	AmtPersonalDeductionMin float64   `db:"amt_personal_deduction_min"`
	CreatedAt               time.Time `db:"created_at"`
	CreatedBy               string    `db:"created_by"`
	UpdatedAt               time.Time `db:"updated_at"`
	UpdatedBy               string    `db:"updated_by"`
}
