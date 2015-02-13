package action

import (
	"fmt"
	"github.com/gohappy/happy"
	"github.com/gohappy/happy/validator"
	"github.com/maxwayt/pkgmanager/plugin"
	"github.com/maxwayt/pkgmanager/storage"
	"github.com/maxwayt/pkgmanager/utility"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type updatePackageAction struct {
	happy.Action
}

func newUpdatePackageAction(context *happy.Context) happy.ActionInterface {

	// Init
	this := &updatePackageAction{
		happy.Action{
			Context: context,
		},
	}

	context.Request.ParseMultipartForm(2048)

	this.AddParameter("name", true, validator.New(validator.Regexp(`^[a-zA-Z0-9-_]+$`), "Invalid Package Name"))
	this.AddParameter("dependencies", false, validator.New(validator.Regexp(`^[a-zA-Z0-9-_]+(,[a-zA-Z0-9-_]+)*$`), "Invalid Dependencies"))
	this.AddParameter("summary", false, validator.New(validator.IsNotEmpty(), "Invalid Summary"))

	return this
}

func (c *updatePackageAction) Run() {

	log.Println("debug:", "name:", c.GetParam("name"))

	file, _, err := c.Context.Request.FormFile("file")
	if err != nil {

		c.AddError(400, "Invalid File")
		return
	}

	defer file.Close()
	mimeType := utility.DetectContentType(file)

	log.Println("debug:", "mimeType:", mimeType)

	// TODO check java file

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

	if deps := c.GetParam("dependencies"); len(deps) != 0 {
		pkg.Dependencies = deps
	}
	if sum := c.GetParam("summary"); len(sum) != 0 {
		pkg.Summary = sum
	}
	pkg.Revision += 1

	filename := fmt.Sprintf(`%s_%d.jar`, c.GetParam("name"), pkg.Revision)

	if err := plugin.Storage.UploadMime(storage.BUCKET_JAR, filename, file, mimeType); err != nil {

		log.Panicln("panic:", err)
	}

	pkg.FileUrl = plugin.Storage.GetFilePublicUrl(storage.BUCKET_JAR, filename)

	if err := col.UpdateId(pkg.Id, pkg); err != nil {

		log.Panicln("panic:", err)
	}

	c.Send(200, `{}`)
}
