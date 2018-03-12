package simulation

import (
	"encoding/json"
	"math"
	"testing"

	"chuckssim.soystudios.com/chuckssim/pkg/bot"
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

func TestEncoding(t *testing.T) {
	expect := `{"bots":[{"id":1,"x":2,"y":3,"a":4.5,"isAutotroph":false}` +
		`,{"id":2,"x":5,"y":6,"a":7.5,"isAutotroph":true}]` +
		`,"type":"state"}`
	sim, err := New()
	if err != nil {
		t.Fatal(err)
	}
	sim.addBot(bot.Bot{X: 2, Y: 3, A: 4.5})
	sim.addBot(bot.Bot{X: 5, Y: 6, A: 7.5, IsAutotroph: true})
	enc, err := json.Marshal(sim)
	if err != nil {
		t.Fatal(err)
	}
	if string(enc) != expect {
		t.Errorf("unexpected encoded\n%s\nexpect\n%s", string(enc), expect)
	}
}
