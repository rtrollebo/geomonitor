package geo

import (
	"testing"
	"time"
)

func TestSearchByTime(t *testing.T) {

	array1 := []GoesXray{}
	i, c := IndexAt(array1, time.Date(
		2025, 01, 01, 12, 03, 00, 0, time.UTC), 1)
	assertTimeWithintRange(t, i, c, 0, 1, len(array1))

	array2 := generateGoesArrayWithPeaks()
	i, c = IndexAt(array2, time.Date( //
		2025, 01, 01, 12, 03, 00, 0, time.UTC), 1)
	assertTimeWithintRange(t, i, c, 6, 1, len(array2))
	i, c = IndexAt(array2, time.Date(
		2025, 01, 01, 12, 01, 00, 0, time.UTC), 1)
	assertTimeWithintRange(t, i, c, 2, 1, len(array2))
	i, c = IndexAt(array2, time.Date(
		2025, 01, 05, 12, 00, 00, 0, time.UTC), 1)
	assertTimeWithintRange(t, i, c, 0, 1, len(array2))
}

func TestSearchByTimeLargeArray(t *testing.T) {
	array1 := generateEquitemporalGoesXrayArray(60)
	i, c := IndexAt(array1, time.Date(
		2025, 01, 01, 12, 15, 00, 0, time.UTC), 2)
	assertTimeWithintRange(t, i, c, 30, 1, len(array1))
	i, c = IndexAt(array1, time.Date(
		2025, 01, 01, 12, 45, 00, 0, time.UTC), 2)
	assertTimeWithintRange(t, i, c, 90, 1, len(array1))
	i, c = IndexAt(array1, time.Date(
		2025, 01, 01, 12, 25, 00, 0, time.UTC), 2)
	assertTimeWithintRange(t, i, c, 50, 1, len(array1))
	i, c = IndexAt(array1, time.Date(
		2025, 01, 01, 12, 00, 00, 0, time.UTC), 2)
	assertTimeWithintRange(t, i, c, 1, 1, len(array1))

	array2 := generateEquitemporalGoesXrayArray(60 * 60 * 24)
	i, c = IndexAt(array2, time.Date(2025, 01, 02, 8, 00, 00, 0, time.UTC), 2)
	assertTimeWithintRange(t, i, c, (12+8)*60*2, 3, len(array2))
}

func TestDetectPeak(t *testing.T) {
	array1 := generateGoesArrayWithPeaks()
	timeFrom := time.Date(2025, 01, 01, 12, 00, 00, 0, time.UTC)
	event, _, err := DetectEvent(array1, timeFrom)
	if err != nil {
		t.Error("Error: ", err)
	}
	if len(event) != 1 || event[0].Class != XRAY_FLARE_M || event[0].Processed == false || event[0].TimeStart.Before(timeFrom) || event[0].TimeEnd.After(timeFrom.Add(4*time.Minute)) {
		t.Error("Failed: ", event)
	}

}

func assertTimeWithintRange(t *testing.T, idx, c, centerValue int, tolerance int, arraysize int) {
	if c > abs((arraysize/2)-centerValue) {
		t.Errorf("Too many iterations")
	}
	if idx < (centerValue-tolerance) || idx > (centerValue+tolerance) {
		t.Errorf("Failed: %d %d", idx, c)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func generateGoesArrayWithPeaks() []GoesXray {
	t1 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 00, 00, 0, time.UTC), Flux: 2.0e-06}
	t2 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 00, 30, 0, time.UTC), Flux: 7.0e-06}
	t3 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 01, 00, 0, time.UTC), Flux: 2.0e-05}
	t4 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 01, 30, 0, time.UTC), Flux: 3.0e-05}
	t5 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 02, 00, 0, time.UTC), Flux: 1.0e-05}
	t6 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 02, 30, 0, time.UTC), Flux: 4.0e-06}
	t7 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 03, 00, 0, time.UTC), Flux: 1.0e-06}
	t8 := GoesXray{TimeTag: time.Date(
		2025, 01, 01, 12, 03, 30, 0, time.UTC), Flux: 3.0e-06}
	array := []GoesXray{t1, t2, t3, t4, t5, t6, t7, t8}
	return array
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
