package simulation

import (
	crand "crypto/rand"
	"encoding/json"
	"io"
	"math"
	"math/big"
	"math/rand"
	"sync/atomic"

	"chuckssim.soystudios.com/chuckssim/pkg/bot"
	errors "github.com/pkg/errors"
)

type (
	Simulation struct {
		RandSource io.Reader `json:"-"`

		BotSize float64 `json:"-"`

		Bots Bots `json:"bots"`

		Type string `json:"type"`

		nextID int64
	}

	Bots map[int64]bot.Bot
)

func (bots Bots) MarshalJSON() ([]byte, error) {
	sl := make([]bot.Bot, len(bots))
	i := 0
	for _, b := range bots {
		sl[i] = b
		i++
	}
	return json.Marshal(sl)
}

func New() (*Simulation, error) {
	sim := &Simulation{}

	sim.Bots = make(map[int64]bot.Bot)
	sim.RandSource = crand.Reader
	sim.BotSize = 10
	sim.Type = "state"

	return sim, nil
}

// NextID returns the next unique ID
func (sim *Simulation) NextID() int64 {
	return atomic.AddInt64(&sim.nextID, 1)
}

// GenerateRandom populates the simulation with a random pupulation
// of autotroph bots.
// They will be contained in the Box 0-boundX, 0-boundY,
// with at least min clusters and at most max clusters.
func (sim *Simulation) GenerateRandom(boundX, boundY float64, min, max int) error {
	seed, err := crand.Int(sim.RandSource, big.NewInt(math.MaxInt64))
	if err != nil {
		return errors.Wrap(err, "error initializing rnd seed")
	}
	rnd := rand.New(rand.NewSource(seed.Int64()))
	maxClusters := rnd.Intn(max+1) + min
	sim.Bots = make(map[int64]bot.Bot, maxClusters*10)
	var rndBot bot.Bot
	for i := 0; i < maxClusters; i++ {
		rndBot = bot.Bot{
			X:           rnd.Float64() * boundX,
			Y:           rnd.Float64() * boundY,
			A:           rnd.Float64() * 2 * math.Pi,
			IsAutotroph: true,
		}
		for {
			sim.addBot(rndBot)
			if rnd.Float64() < 0.3 {
				break
			}
			a := rnd.Float64() * 2 * math.Pi
			rndBot.X, rndBot.Y = sim.placeNextTo(rndBot.X, rndBot.Y, a)
			rndBot.ID = sim.NextID()
		}
	}
	return nil
}

func (sim *Simulation) addBot(b bot.Bot) {
	b.ID = sim.NextID()
	sim.Bots[b.ID] = b
}

func (sim *Simulation) placeNextTo(x, y, a float64) (newX, newY float64) {
	return x + sim.BotSize*math.Cos(a), y + sim.BotSize*math.Sin(a)
}
