package utxo

type (
	Entry struct {
		UserId string
		Amount  float32
	}

	UTXO map[string]float32
)
