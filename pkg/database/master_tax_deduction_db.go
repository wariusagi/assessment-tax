package database

import (
	"fmt"
)

var CacheMasterTaxDeduction = make(map[int]MasterTaxDeduction)

func (r repositoryDB) GetMasterTaxDeduction(cycleYear int) (MasterTaxDeduction, error) {
	cache, ok := CacheMasterTaxDeduction[cycleYear]
	if ok {
		return cache, nil
	}

	stmt, err := r.db.Prepare(
		`
		SELECT cycle_year, amt_k_receipt_max, amt_donation_max, amt_personal_deduction_min
		FROM master_tax_deduction
		WHERE cycle_year = $1
		`,
	)
	if err != nil {
		return MasterTaxDeduction{}, fmt.Errorf(`can't prepare statment: %v`, err)
	}

	row := stmt.QueryRow(cycleYear)
	result := MasterTaxDeduction{}

	err = row.Scan(&result.CycleYear, &result.AmtKReceiptMax, &result.AmtDonationMax, &result.AmtPersonalDeductionMin)
	if err != nil {
		return MasterTaxDeduction{}, fmt.Errorf(`can't scan row into variables: %v`, err)
	}
	CacheMasterTaxDeduction[cycleYear] = result
	return result, nil
}
