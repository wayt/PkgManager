package action

import (
	"github.com/gohappy/happy"
	"github.com/gohappy/happy/validator"
	"github.com/maxwayt/pkgmanager/plugin"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type deletePackageAction struct {
	happy.Action
}

func newDeletePackageAction(context *happy.Context) happy.ActionInterface {

	// Init
	this := &deletePackageAction{
		happy.Action{
			Context: context,
		},
	}

	this.AddParameter("name", true, validator.New(validator.Regexp(`^[a-zA-Z0-9-_]+$`), "Invalid Package Name"))

	return this
}

func (c *deletePackageAction) Run() {

	log.Println("debug:", "name:", c.GetParam("name"))

	col := plugin.MongoDB.C("packages")

	res := col.Find(bson.M{"name": c.GetParam("name")})

	count, err := res.Count()
	if err != nil {

		log.Panicln("panic:", err)
	}

	if count == 0 {

		c.AddError(404, "Package not Found")
		return
	}

	var pkg packageInfo
	if err := res.One(&pkg); err != nil {

		log.Panicln("panic:", err)
	}

	if err := col.RemoveId(pkg.Id); err != nil {

		log.Panicln("panic:", err)
	}

	c.Send(200, `{}`)
}
