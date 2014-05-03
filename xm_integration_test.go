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
		"domain":  os.Getenv("XM_DOMAIN"),  // example.org,
		"ml":      os.Getenv("XM_ML"),      // mailing-list
	}
}

func setupIntegrationClient() (*Client, map[string]string) {
	vars := getIntegrationTestVars()
	cl, err := Dial(vars["network"], vars["host"])
	if err != nil {
		panic(err)
	}
	if err := cl.Authenticate(vars["user"], vars["pass"]); err != nil {
		panic(err)
	}
	return cl, vars
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

	if err := cl.Noop(); err != nil {
		t.Fatal(err)
	}
}

func TestMailingListUserMod(t *testing.T) {
	cl, vars := setupIntegrationClient()
	defer cl.Close()
	if err := cl.MailingListAddUser(vars["domain"], vars["ml"], "test@example.org", ""); err != nil {
		t.Fatal(err)
	}
	if err := cl.MailingListDeleteUser(vars["domain"], vars["ml"], "test@example.org"); err != nil {
		t.Fatal(err)
	}
}
