package rankings

import (
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
		 _, err := ConnectToPSQL()
		 if err != nil {
		 		t.Errorf("err = %v; want <nil>", err)
		 }
}