package initializers

import (
	"github.com/google/wire"

	"github.com/F0urward/proftwist-backend/internal/wire/sets"
)

var HTTPServerSet = wire.NewSet(
	sets.CommonSet,
	sets.RoadmapInfoSet,
	sets.RoadmapSet,
	sets.AuthSet,
)
