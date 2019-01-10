package chatapigateway

import (
	"testing"
)

func TestLive(t *testing.T) {
	c := &Connection{ConnectionID: "123", Channel: "abc"}
	sess, err := NewAwsSession("us-west-1", "jds")
	if err != nil {
		t.Error(err)
	}

	err = Put(sess, *c)
	if err != nil {
		t.Error(err)
	}

	rcget, err := Get(sess, "abc", "123")
	if err != nil {
		t.Error(err)
	}
	if rcget.Channel != "abc" || rcget.ConnectionID != "123" {
		t.Errorf("get err %v", rcget)
	}

	rc, err := GetChannelConnections(sess, "abc")
	if err != nil {
		t.Error(err)
	}
	if len(rc) < 1 {
		t.Error("no results")
	}

	err = Delete(sess, "abc", "123")
	if err != nil {
		t.Error(err)
	}
}
