package app

import (
	"github.com/ahnaftahmid39/gator/internal/config"
	"github.com/ahnaftahmid39/gator/internal/database"
)

type State struct {
	Cfg *config.Config
	Db  *database.Queries
}
