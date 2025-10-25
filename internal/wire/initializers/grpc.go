package initializers

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/internal/wire/sets"
)

var GRPCServerSet = wire.NewSet(
	sets.CommonSet,
	sets.RoadmapSet,
	sets.RoadmapInfoSet,
)
