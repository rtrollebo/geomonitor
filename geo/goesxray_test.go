package geo

import (
	"fmt"
	"testing"
	"time"
)

func TestSearchByTime(t *testing.T) {

	t1 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 00, 00, 0, time.UTC)}
	t2 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 00, 30, 0, time.UTC)}
	t3 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 01, 00, 0, time.UTC)}
	t4 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 01, 30, 0, time.UTC)}
	t5 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 02, 00, 0, time.UTC)}
	t6 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 02, 30, 0, time.UTC)}
	t7 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 03, 00, 0, time.UTC)}
	t8 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 03, 30, 0, time.UTC)}

	array := []GoesXray{t1, t2, t3, t4, t5, t6, t7, t8}

	idx1, c1 := SearchByTime(array, time.Date(
		2025, 01, 01, 12, 03, 00, 0, time.UTC))
	if c1 > 5 {
		t.Errorf("Too many iterations")
	}
	if idx1 < 5 && idx1 > 7 {
		t.Errorf(fmt.Sprintf("Failed: idx1=%d, %d", idx1, c1))
	}

	idx2, c2 := SearchByTime(array, time.Date(
		2025, 01, 01, 12, 01, 00, 0, time.UTC))
	if c2 > 5 {
		t.Errorf("Too many iterations")
	}
	if idx2 < 1 && idx2 > 3 {
		t.Errorf(fmt.Sprintf("Failed: idx2=%d, %d", idx2, c2))
	}

}
