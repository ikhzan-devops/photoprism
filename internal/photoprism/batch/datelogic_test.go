package batch

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestComputeDateChange_LeapClamp(t *testing.T) {
	base := time.Date(2024, 2, 29, 11, 29, 54, 0, time.UTC)
	newLocal, y, m, d := ComputeDateChange(base, 2024, 2, 29, ActionNone, 0, ActionNone, 0, ActionUpdate, 2025)

	assert.Equal(t, 2025, y)
	assert.Equal(t, 2, m)
	assert.Equal(t, 28, d)
	assert.Equal(t, 2025, newLocal.Year())
	assert.Equal(t, time.February, newLocal.Month())
	assert.Equal(t, 28, newLocal.Day())
}

func TestComputeDateChange_DayUnknown_YearNotForcedUnknown(t *testing.T) {
	base := time.Date(2023, 11, 12, 9, 0, 0, 0, time.UTC)
	newLocal, y, m, d := ComputeDateChange(base, 2023, 11, 12, ActionUpdate, -1, ActionNone, 0, ActionUpdate, 2000)

	assert.Equal(t, 2000, y)
	assert.Equal(t, 11, m)
	assert.Equal(t, -1, d)
	assert.Equal(t, 1, newLocal.Day()) // day used for TakenAtLocal when unknown
}

func TestComputeDateChange_MonthUnknown_ForcesYearUnknown(t *testing.T) {
	base := time.Date(2020, 4, 30, 8, 0, 0, 0, time.UTC)
	newLocal, y, m, d := ComputeDateChange(base, 2020, 4, 30, ActionNone, 0, ActionUpdate, -1, ActionUpdate, 2000)

	assert.Equal(t, -1, y)
	assert.Equal(t, -1, m)
	assert.Equal(t, 30, d) // day reflects base month clamped result
	assert.Equal(t, 30, newLocal.Day())
	assert.Equal(t, time.April, newLocal.Month())
}

func TestComputeDateChange_MixedMonth_UpdateDayClampsPerPhoto(t *testing.T) {
	baseMar := time.Date(2020, 3, 20, 11, 29, 54, 0, time.UTC)
	newLocal1, y1, m1, d1 := ComputeDateChange(baseMar, 2020, 3, 20, ActionUpdate, 31, ActionNone, 0, ActionNone, 0)
	assert.Equal(t, 2020, y1)
	assert.Equal(t, 3, m1)
	assert.Equal(t, 31, d1)
	assert.Equal(t, 31, newLocal1.Day())

	baseApr := time.Date(2020, 4, 20, 11, 29, 54, 0, time.UTC)
	newLocal2, y2, m2, d2 := ComputeDateChange(baseApr, 2020, 4, 20, ActionUpdate, 31, ActionNone, 0, ActionNone, 0)
	assert.Equal(t, 2020, y2)
	assert.Equal(t, 4, m2)
	assert.Equal(t, 30, d2)
	assert.Equal(t, 30, newLocal2.Day())
}

// Case: Current Month is Unknown (-1), update Day=31 → keep Month unknown, clamp per base month
func TestComputeDateChange_UnknownCurrentMonth_DayUpdateClampsPerPhoto(t *testing.T) {
	// Base in March, current month unknown
	baseMar := time.Date(2020, 3, 20, 11, 29, 54, 0, time.UTC)
	newLocal1, y1, m1, d1 := ComputeDateChange(baseMar, 2020, -1, 20, ActionUpdate, 31, ActionNone, 0, ActionNone, 0)
	assert.Equal(t, 2020, y1)
	assert.Equal(t, -1, m1) // keep unknown
	assert.Equal(t, 31, d1)
	assert.Equal(t, time.March, newLocal1.Month())
	assert.Equal(t, 31, newLocal1.Day())

	// Base in April, current month unknown
	baseApr := time.Date(2020, 4, 20, 11, 29, 54, 0, time.UTC)
	newLocal2, y2, m2, d2 := ComputeDateChange(baseApr, 2020, -1, 20, ActionUpdate, 31, ActionNone, 0, ActionNone, 0)
	assert.Equal(t, 2020, y2)
	assert.Equal(t, -1, m2) // keep unknown
	assert.Equal(t, 30, d2) // clamp to Apr 30
	assert.Equal(t, time.April, newLocal2.Month())
	assert.Equal(t, 30, newLocal2.Day())
}

func TestComputeDateChange_MixedDay_UpdateMonthToFeb2020(t *testing.T) {
	base1 := time.Date(2020, 4, 20, 11, 29, 54, 0, time.UTC)
	newLocal1, y1, m1, d1 := ComputeDateChange(base1, 2020, 4, 20, ActionNone, 0, ActionUpdate, 2, ActionNone, 0)
	assert.Equal(t, 2020, y1)
	assert.Equal(t, 2, m1)
	assert.Equal(t, 20, d1)
	assert.Equal(t, time.February, newLocal1.Month())
	assert.Equal(t, 20, newLocal1.Day())

	base2 := time.Date(2020, 3, 31, 11, 29, 54, 0, time.UTC)
	newLocal2, y2, m2, d2 := ComputeDateChange(base2, 2020, 3, 31, ActionNone, 0, ActionUpdate, 2, ActionNone, 0)
	assert.Equal(t, 2020, y2)
	assert.Equal(t, 2, m2)
	assert.Equal(t, 29, d2)
	assert.Equal(t, time.February, newLocal2.Month())
	assert.Equal(t, 29, newLocal2.Day())
}

func TestComputeDateChange_Mixed_UpdateYearTo2021(t *testing.T) {
	base1 := time.Date(2020, 2, 29, 11, 29, 54, 0, time.UTC)
	newLocal1, y1, m1, d1 := ComputeDateChange(base1, 2020, 2, 29, ActionNone, 0, ActionNone, 0, ActionUpdate, 2021)
	assert.Equal(t, 2021, y1)
	assert.Equal(t, 2, m1)
	assert.Equal(t, 28, d1)
	assert.Equal(t, 2021, newLocal1.Year())
	assert.Equal(t, 28, newLocal1.Day())

	base2 := time.Date(2020, 3, 31, 11, 29, 54, 0, time.UTC)
	newLocal2, y2, m2, d2 := ComputeDateChange(base2, 2020, 3, 31, ActionNone, 0, ActionNone, 0, ActionUpdate, 2021)
	assert.Equal(t, 2021, y2)
	assert.Equal(t, 3, m2)
	assert.Equal(t, 31, d2)
	assert.Equal(t, 2021, newLocal2.Year())
	assert.Equal(t, 31, newLocal2.Day())
}
