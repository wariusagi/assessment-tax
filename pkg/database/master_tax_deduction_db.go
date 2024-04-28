package database

import (
	"fmt"
	"log"
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
		return MasterTaxDeduction{}, fmt.Errorf(`prepare statement failed: %v`, err)
	}

	row := stmt.QueryRow(cycleYear)
	result := MasterTaxDeduction{}

	err = row.Scan(&result.CycleYear, &result.AmtKReceiptMax, &result.AmtDonationMax, &result.AmtPersonalDeductionMin)
	if err != nil {
		return MasterTaxDeduction{}, fmt.Errorf(`scan row into variables failed: %v`, err)
	}
	CacheMasterTaxDeduction[cycleYear] = result
	return result, nil
}

func (r repositoryDB) UpdateAmtPersonalDeductionDeduction(cycleYear int, amtPersonalDeduction float64) error {
	_, err := r.db.Exec(`
		UPDATE master_tax_deduction
		SET amt_personal_deduction_min = $1, updated_at = NOW(), updated_by = $2
		WHERE cycle_year = $3
	`, amtPersonalDeduction, "ADMIN", cycleYear)

	if err != nil {
		return fmt.Errorf("update amt_personal_deduction_min to table master_tax_deduction failed: %v", err)
	}

	delete(CacheMasterTaxDeduction, cycleYear)

	log.Printf("amt_personal_deduction_min updated [%v][%v]", cycleYear, amtPersonalDeduction)
	return nil
}

func (r repositoryDB) UpdateAmtKReceiptDeduction(cycleYear int, amtKReceipt float64) error {
	_, err := r.db.Exec(`
		UPDATE master_tax_deduction
		SET amt_k_receipt_max = $1, updated_at = NOW(), updated_by = $2
		WHERE cycle_year = $3
	`, amtKReceipt, "ADMIN", cycleYear)

	if err != nil {
		return fmt.Errorf("update amt_k_receipt_max to table master_tax_deduction failed: %v", err)
	}

	delete(CacheMasterTaxDeduction, cycleYear)

	log.Printf("amt_k_receipt_max updated [%v][%v]", cycleYear, amtKReceipt)
	return nil
}
