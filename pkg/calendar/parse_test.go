package calendar

import (
	"testing"
)

func TestDirSrvStartTime(t *testing.T) {
	dates, err := ProcessInput(`
01.01.12
9-6-1988
31/12/2023
1/1/2015`)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("dates: %v", dates)
	t.Logf("len: %d", len(dates))
}
