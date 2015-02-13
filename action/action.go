package action

import (
	"github.com/gohappy/happy"
)

func RegisterActions(app *happy.API) {

	app.AddRoute("POST", "/package", newCreatePackageAction)
	app.AddRoute("PUT", "/package", newUpdatePackageAction)
	app.AddRoute("GET", "/package/:name", newGetPackageAction)
	app.AddRoute("DELETE", "/package/:name", newDeletePackageAction)

	app.AddRoute("GET", "/packages/:name", newSearchPackageAction)
}
