package bot

type (
	Bot struct {
		X, Y float64
		A    float64

		Energy  float64
		Signals Signals
	}

	Signals struct {
		A, B, C, D, E, F, G, H, I, J, K, L, M, N, O, P, Q, R, S, T, U, V, W, X, Y, Z float64
	}
)
