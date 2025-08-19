package types

type (
	Account struct {
		ID       int
		Username string
		Balance  int
		Depth    int
	}

	Log struct {
		LogType  string
		Date     string
		Username string
	}
)
