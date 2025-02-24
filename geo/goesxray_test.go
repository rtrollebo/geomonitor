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
	assertTimeWithinRange(t, array, time.Date(
		2025, 01, 01, 12, 03, 00, 0, time.UTC), 6)
	assertTimeWithinRange(t, array, time.Date(
		2025, 01, 01, 12, 01, 00, 0, time.UTC), 2)
}

func assertTimeWithinRange(t *testing.T, array []GoesXray, timeValue time.Time, centerValue int) {
	idx, c := SearchByTime(array, timeValue)
	if c > (len(array) / 2) {
		t.Errorf("Too many iterations")
	}
	if idx < (centerValue-1) || idx > (centerValue+1) {
		t.Errorf(fmt.Sprintf("Failed: idx=%d, %d", idx, c))
	}

}
