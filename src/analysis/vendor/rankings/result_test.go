package rankings

import (
	"testing"
	// "database/sql"
	// _ "github.com/lib/pq"
)

func TestGetTime(t *testing.T) {
		 time := GetTime("1:01:01")
		 if time != (61*60 + 1) {
			t.Errorf("t = %v; want 3661", time)
		 }
}

func TestGetTime2(t *testing.T) {
	time := GetTime("1:01:01.5")
	if time != (61.0*60.0 + 1.5) {
	 t.Errorf("t = %v; want 3661.5", time)
	}
}
