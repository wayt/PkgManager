package action

import (
	"encoding/json"
	"github.com/gohappy/happy"
	"github.com/gohappy/happy/validator"
	"github.com/wayt/pkgmanager/plugin"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type getPackageAction struct {
	happy.Action
}

func newGetPackageAction(context *happy.Context) happy.ActionInterface {

	// Init
	this := &getPackageAction{
		happy.Action{
			Context: context,
		},
	}

	this.AddParameter("name", true, validator.New(validator.Regexp(`^[a-zA-Z0-9-_]+$`), "Invalid Package Name"))

	return this
}

func (c *getPackageAction) Run() {

	col := plugin.MongoDB.C("packages")

	res := col.Find(bson.M{"name": c.GetParam("name")})

	cnt, err := res.Count()
	if err != nil {

		log.Panicln("panic:", err)
	}
	if cnt == 0 {

		c.AddError(404, "Package Not Found")
		return
	}

	var pkg packageInfo
	if err := res.One(&pkg); err != nil {

		log.Panicln("panic:", err)
	}

	data, err := json.Marshal(pkg)
	if err != nil {

		log.Panicln("panic:", err)
	}

	c.Send(200, string(data))
}
