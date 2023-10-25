package calendar

import (
	"testing"
)

func TestDirSrvStartTime(t *testing.T) {
	_, err := ProcessInput(`
	11/05/2023
	10/07/2023
	15/07/2023
	11/08/2023
	5/11/2023
`)
	if err != nil {
		t.Error(err)
		return
	}
}
