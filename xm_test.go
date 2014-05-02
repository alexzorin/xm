package xm

import (
	"testing"
)

func TestParse(t *testing.T) {
	// Easy server banner test
	code, msg, err := parseLine("-00000 <1399064990.140310407509760@1.2.3.4> XMail 1.27 CTRL Server; Sat, 3 May 2014 07:09:50 +1000")
	if code != 0 || msg != "<1399064990.140310407509760@1.2.3.4> XMail 1.27 CTRL Server; Sat, 3 May 2014 07:09:50 +1000" || err != nil {
		t.Fatal("Couldn't parse server banner", err)
	}

	// Syntax error test
	code, msg, err = parseLine("-00103 Bad CTRL command syntax")
	if code != 103 || msg != "Bad CTRL command syntax" || err != nil {
		t.Fatal("Couldn't parse syntax error", err)
	}

	// OK test
	code, msg, err = parseLine("+00100 OK")
	if code != 100 || msg != "OK" || err != nil {
		t.Fatal("Couldn't parse OK", code, msg, err)
	}
}
