package sender

import "btc_service/src/model"

type (
	sender   struct{}
	database interface {
		Save()
		Read() model.Email
	}
)
