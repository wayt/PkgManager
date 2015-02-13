package action

import (
	"fmt"
	"github.com/gohappy/happy"
	"github.com/gohappy/happy/validator"
	"github.com/maxwayt/pkgmanager/plugin"
	"github.com/maxwayt/pkgmanager/storage"
	"github.com/maxwayt/pkgmanager/utility"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"strings"
)

type createPackageAction struct {
	happy.Action
}

type packageInfo struct {
	Id           bson.ObjectId `bson:"_id" json:"-"`
	Name         string        `bson:"name" json:"name"`
	AuthorEmail  string        `bson:"author_email" json:"author_email"`
	Dependencies string        `bson:"dependencies" json:"dependencies,omitempty"`
	Revision     int64         `bson:"revision" json:"revision"`
	Summary      string        `bson:"summary" json:"summary,omitempty"`
	FileUrl      string        `bson:"file_url" json:"file_url"`
}

func newCreatePackageAction(context *happy.Context) happy.ActionInterface {

	// Init
	this := &createPackageAction{
		happy.Action{
			Context: context,
		},
	}

	context.Request.ParseMultipartForm(2048)

	this.AddParameter("name", true, validator.New(validator.Regexp(`^[a-zA-Z0-9-_]+$`), "Invalid Package Name"))
	this.AddParameter("author_email", true, validator.New(validator.IsEmail(), "Invalid Author Email"))
	this.AddParameter("dependencies", false, validator.New(validator.Regexp(`^[a-zA-Z0-9-_]+(,[a-zA-Z0-9-_]+)*$`), "Invalid Dependencies"))
	this.AddParameter("summary", false, validator.New(validator.IsNotEmpty(), "Invalid Summary"))

	return this
}

func (c *createPackageAction) Run() {

	log.Println("debug:", "name:", c.GetParam("name"))

	contentType := c.Context.Request.Header.Get("Content-Type")
	headers := strings.Split(contentType, `;`)
	if headers[0] != "multipart/form-data" {

		log.Println("debug:", "contentType:", headers[0])
		c.AddError(400, "Invalid `Content-Type` Header, need `multipart/form-data`")
		return
	}

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

	if count > 0 {

		c.AddError(409, "Duplicate Package Name")
		return
	}

	filename := fmt.Sprintf(`%s_%d.jar`, c.GetParam("name"), 1)

	if err := plugin.Storage.UploadMime(storage.BUCKET_JAR, filename, file, mimeType); err != nil {

		log.Panicln("panic:", err)
	}

	if err := col.Insert(&packageInfo{
		Id:           bson.NewObjectId(),
		Name:         c.GetParam("name"),
		AuthorEmail:  c.GetParam("author_email"),
		Dependencies: c.GetParam("dependencies"),
		Summary:      c.GetParam("summary"),
		Revision:     1,
		FileUrl:      plugin.Storage.GetFilePublicUrl(storage.BUCKET_JAR, filename),
	}); err != nil {

		log.Panicln("panic:", err)
	}

	if err := col.EnsureIndex(mgo.Index{
		Key:        []string{"name"},
		Unique:     true,
		DropDups:   true,
		Background: true, // See notes.
		Sparse:     true,
	}); err != nil {
		log.Panicln("panic:", err)
	}

	c.Send(200, `{}`)
}
