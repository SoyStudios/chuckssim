package simulation

import (
	"math"
	"testing"
)

type placementTest struct {
	X       float64
	Y       float64
	A       float64
	expectX float64
	expectY float64
}

func providePlacementTests() []placementTest {
	return []placementTest{
		{0, 0, 0, 5, 0},
		{0, 0, 0.5 * math.Pi, 0, 5},
		{0, 0, math.Pi, -5, 0},
		{0, 0, 1.5 * math.Pi, 0, -5},
		{0, 0, 2 * math.Pi, 5, 0},
	}
}

func TestPlacement(t *testing.T) {
	sim, err := New()
	if err != nil {
		t.Fatal(err)
	}
	sim.BotSize = 5
	for _, pt := range providePlacementTests() {
		x, y := sim.placeNextTo(pt.X, pt.Y, pt.A)
		if math.Abs(pt.expectX-x) > 0.0001 {
			t.Errorf("unexpected X %.2f, expect %.2f", x, pt.expectX)
			return
		}
		if math.Abs(pt.expectY-y) > 0.0001 {
			t.Errorf("unexpected Y %.2f, expect %.2f", y, pt.expectY)
		}
	}
}
