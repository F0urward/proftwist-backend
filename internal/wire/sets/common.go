package sets

import (
	"github.com/google/wire"

	db "github.com/F0urward/proftwist-backend/internal/infrastructure/db/postgres"
)

var CommonSet = wire.NewSet(
	db.New,
)
