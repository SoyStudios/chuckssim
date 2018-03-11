package bot

import "encoding/json"

type (
	botDisplay struct {
		ID int64 `json:"id"`
		// x, y and angle
		X           float64 `json:"x"`
		Y           float64 `json:"y"`
		A           float64 `json:"a"`
		IsAutotroph bool    `json:"isAutotroph"`
	}

	Bot struct {
		ID int64 `json:"id"`
		// x, y and angle
		X float64 `json:"x"`
		Y float64 `json:"y"`
		A float64 `json:"a"`
		// delta V
		DV float64

		Energy float64

		IsAutotroph bool

		Signals Signals
	}

	Signal [9]float64

	Signals struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z Signal
	}
)

func (b Bot) display() botDisplay {
	return botDisplay{
		ID:          b.ID,
		X:           b.X,
		Y:           b.Y,
		A:           b.A,
		IsAutotroph: b.IsAutotroph,
	}
}

func (b Bot) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.display())
}
