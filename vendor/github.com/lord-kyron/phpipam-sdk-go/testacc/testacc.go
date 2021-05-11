// Package testacc contains helper methods for running acceptance tests.
package testacc

import (
	"os"
	"testing"
)

// SkipIfNotAcc is designed to skip an integration test if TESTACC is not set.
func SkipIfNotAcc(t *testing.T) {
	if os.Getenv("TESTACC") == "" {
		t.Skipf("Skipping integration test as TESTACC is not set.")
	}
}

// SkipIfCustomNested is designed to skip an integration test if
// TESTACC_CUSTOM_NESTED is set.
func SkipIfCustomNested(t *testing.T) {
	if os.Getenv("TESTACC_CUSTOM_NESTED") != "" {
		t.Skipf("Skipping non-nested custom field test because TESTACC_CUSTOM_NESTED is set")
	}
}

// PanicIfMissingEnv is designed to panic if the following environment variables
// are not set:
//
//  * PHPIPAM_APP_ID
//  * PHPIPAM_ENDPOINT_ADDR
//  * PHPIPAM_PASSWORD
//  * PHPIPAM_USER_NAME
//
// Acceptance tests cannot continue if these are not set so there is no point
// in continuing.
func PanicIfMissingEnv() {
	if os.Getenv("PHPIPAM_APP_ID") == "" || os.Getenv("PHPIPAM_ENDPOINT_ADDR") == "" || os.Getenv("PHPIPAM_PASSWORD") == "" || os.Getenv("PHPIPAM_USER_NAME") == "" {
		panic("Please ensure the environment variables PHPIPAM_APP_ID, PHPIPAM_ENDPOINT_ADDR, PHPIPAM_PASSWORD, and PHPIPAM_USER_NAME are set for acceptance tests")
	}
}

// VetAccConditions is a meta-function that ensures that an acceptance test
// meets the conditions necessary to continue.
func VetAccConditions(t *testing.T) {
	SkipIfNotAcc(t)
	PanicIfMissingEnv()
}
