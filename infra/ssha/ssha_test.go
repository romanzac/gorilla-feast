package ssha

import "testing"

func TestValidatePassword(t *testing.T) {

	// OK plain text and OK password from DB
	pwd := "roman123"
	pwdDB := "{SSHA}bQ/+rmGtWkhGbegZvRXch3nJ9puRo63xU/kFq5pMtmhkrgDolN6RkZBQXfL69rS8SPtFrw=="
	res, err := ValidatePassword(pwd, pwdDB)
	if res == false || err != nil {
		t.Errorf("Failed validation for correct plaintext and DB pair")
	}

	// Wrong plaintext, ok DB
	pwd = "rfsaftfoman123"
	pwdDB = "{SSHA}bQ/+rmGtWkhGbegZvRXch3nJ9puRo63xU/kFq5pMtmhkrgDolN6RkZBQXfL69rS8SPtFrw=="
	res, err = ValidatePassword(pwd, pwdDB)
	if res == true {
		t.Errorf("Failed validation for wrong plaintext, ok DB pair")
	}

	// Ok plain text, wrong DB
	pwd = "roman123"
	pwdDB = "{SSHA}bQ/+rmGtWkhGbegZvRXch3nJ9puRo63xU/kFq5pMtmhkrgDolN6RkZBQXfL693S8SPtFrw=="
	res, err = ValidatePassword(pwd, pwdDB)
	if res == true {
		t.Errorf("Failed validation for OK plain text, wrong DB pair")
	}

	// Missing {SSHA} prefix in DB, OK plain text and OK password from DB
	pwd = "roman123"
	pwdDB = "bQ/+rmGtWkhGbegZvRXch3nJ9puRo63xU/kFq5pMtmhkrgDolN6RkZBQXfL69rS8SPtFrw=="
	res, err = ValidatePassword(pwd, pwdDB)
	if err == nil {
		t.Errorf("Failed validation for correct plaintext and DB pair")
	}
}
