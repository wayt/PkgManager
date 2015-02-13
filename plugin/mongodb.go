package plugin

import (
	"gopkg.in/mgo.v2"
)

var MongoDB *mgo.Database

func SetupMongoDB() error {

	mongo, err := mgo.Dial(Config.MongoDB.Servers)
	if err != nil {
		return err
	}
	//defer mongo.Close()

	MongoDB = mongo.DB(Config.MongoDB.Db)

	return nil
}
