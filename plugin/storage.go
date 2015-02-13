package plugin

import (
	"github.com/maxwayt/pkgmanager/storage"
)

var Storage *storage.Storage

func SetupStorage() {

	Storage, _ = storage.NewStorage(Config.Storage)

}
