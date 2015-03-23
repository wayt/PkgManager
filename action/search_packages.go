package action

import (
	"encoding/json"
	"github.com/gohappy/happy"
	"github.com/gohappy/happy/validator"
	"github.com/wayt/pkgmanager/plugin"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type searchPackageAction struct {
	happy.Action
}

func newSearchPackageAction(context *happy.Context) happy.ActionInterface {

	// Init
	this := &searchPackageAction{
		happy.Action{
			Context: context,
		},
	}

	this.AddParameter("name", true, validator.New(validator.Regexp(`^[a-zA-Z0-9-_]+$`), "Invalid Package Name"))

	return this
}

func (c *searchPackageAction) Run() {

	col := plugin.MongoDB.C("packages")

	res := col.Find(bson.M{
		"name": bson.M{
			"$regex": bson.RegEx{
				`^.*` + c.GetParam("name") + `.*$`,
				"i",
			},
		},
	})

	cnt, err := res.Count()
	if err != nil {

		log.Panicln("panic:", err)
	}
	if cnt == 0 {

		c.Send(200, `[]`)
		return
	}

	pkgs := []packageInfo{}

	iter := res.Iter()
	if err := iter.All(&pkgs); err != nil {

		log.Panicln("panic:", err)
	}

	if err := iter.Close(); err != nil {
		log.Panicln("panic:", err)
	}

	data, err := json.Marshal(pkgs)
	if err != nil {

		log.Panicln("panic:", err)
	}

	c.Send(200, string(data))
}
