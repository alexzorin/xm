// +build integration

package xm

import (
	"os"
	"testing"
)

func getIntegrationTestVars() map[string]string {
	return map[string]string{
		"host":    os.Getenv("XM_HOST"),    // localhost:6017
		"network": os.Getenv("XM_NETWORK"), // tcp
		"user":    os.Getenv("XM_USER"),    // admin
		"pass":    os.Getenv("XM_PASS"),    // adminpassword
	}
}

func TestConnect(t *testing.T) {
	vars := getIntegrationTestVars()
	cl, err := Dial(vars["network"], vars["host"])
	if err != nil {
		t.Fatal(err)
	}
	defer cl.Close()

	if err := cl.Authenticate(vars["user"], vars["pass"]); err != nil {
		t.Fatal(err)
	}
}
