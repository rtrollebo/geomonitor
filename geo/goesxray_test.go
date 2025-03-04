package geo

import (
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

	array1 := []GoesXray{}
	i, c := IndexAt(array1, time.Date(
		2025, 01, 01, 12, 03, 00, 0, time.UTC), 1)
	assertTimeWithintRange(t, i, c, 0)

	array2 := []GoesXray{t1, t2, t3, t4, t5, t6, t7, t8}
	i, c = IndexAt(array2, time.Date( //
		2025, 01, 01, 12, 03, 00, 0, time.UTC), 1)
	assertTimeWithintRange(t, i, c, 6)
	i, c = IndexAt(array2, time.Date(
		2025, 01, 01, 12, 01, 00, 0, time.UTC), 1)
	assertTimeWithintRange(t, i, c, 2)
	i, c = IndexAt(array2, time.Date(
		2025, 01, 05, 12, 00, 00, 0, time.UTC), 1)
	assertTimeWithintRange(t, i, c, 0)
}

func TestSearchByTimeLargeArray(t *testing.T) {
	array := generateEquitemporalGoesXrayArray(60)
	i, c := IndexAt(array, time.Date(
		2025, 01, 01, 12, 15, 00, 0, time.UTC), 2)
	assertTimeWithintRange(t, i, c, 30)
	i, c = IndexAt(array, time.Date(
		2025, 01, 01, 12, 45, 00, 0, time.UTC), 2)
	assertTimeWithintRange(t, i, c, 90)
	i, c = IndexAt(array, time.Date(
		2025, 01, 01, 12, 25, 00, 0, time.UTC), 2)
	assertTimeWithintRange(t, i, c, 50)
}

func assertTimeWithintRange(t *testing.T, idx, c, centerValue int) {
	if c > (centerValue / 2) {
		t.Errorf("Too many iterations")
	}
	if idx < (centerValue-1) || idx > (centerValue+1) {
		t.Errorf("Failed: %d %d", idx, c)
	}
}

func generateEquitemporalGoesXrayArray(minutes int) []GoesXray {
	initialDate := time.Date(2025, 01, 01, 12, 00, 00, 0, time.UTC)
	currentDate := initialDate
	n := minutes * 2
	array := make([]GoesXray, n)
	for i := 0; i < n; i++ {
		if i%2 == 0 {
			currentDate = initialDate.Add(time.Duration(i/2) * time.Minute)
		}
		array[i] = GoesXray{TimeTag: currentDate}
	}
	return array
}
