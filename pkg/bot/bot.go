package bot

type (
	Bot struct {
		ID int64
		// x, y and angle
		X, Y float64
		A    float64
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
