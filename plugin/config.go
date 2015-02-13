package plugin

import (
	"encoding/json"
	"github.com/maxwayt/pkgmanager/storage"
	"log"
	"os"
	"runtime"
)

var Config struct {
	BindIp     string `json:"bind_ip"`
	BindPort   int    `json:"bind_port"`
	LogFile    string `json:"log_file"`
	GoMaxProcs int    `json:"go_max_procs"`
	MongoDB    struct {
		Servers string `json:"servers"`
		Db      string `json:"db"`
	} `json:"mongodb"`
	TempDir string                `json:"temp_dir"`
	Storage storage.StorageConfig `json:"storage"`
}

func ReadConfig(filename string) error {

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&Config); err != nil {
		return err
	}

	if Config.LogFile != "" {

		f, err := os.OpenFile(Config.LogFile, os.O_RDONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		log.SetOutput(f)
	}

	if Config.GoMaxProcs > 0 {
		runtime.GOMAXPROCS(Config.GoMaxProcs)
	}

	return nil
}
