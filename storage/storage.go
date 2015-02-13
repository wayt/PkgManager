package storage

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type StorageConfig struct {
	Host       string `json:"host"`
	User       string `json:"user"`
	PublicHost string `json:"public_host"`
}

type Bucket string

const (
	BUCKET_JAR = "jar"
)

type Storage struct {
	Client *http.Client

	Config StorageConfig
}

func NewStorage(config StorageConfig) (*Storage, error) {

	this := &Storage{}

	this.Client = &http.Client{}
	this.Config = config

	return this, nil
}

func DetectContentType(file io.ReadSeeker) string {

	// first 512 bytes are used to evaluate mime type
	first512 := make([]byte, 512)
	file.Read(first512)
	file.Seek(0, 0)
	return http.DetectContentType(first512)
}

func (this *Storage) GetFilePublicUrl(bucket Bucket, filename string) string {

	return fmt.Sprintf("%s/%s/%s/%s", this.Config.PublicHost, this.Config.User, bucket, filename)
}

func (this *Storage) UploadMime(bucket Bucket, objectname string, contents io.Reader, mimeType string) error {

	url := fmt.Sprintf("%s/%s/%s/%s", this.Config.Host, this.Config.User, bucket, objectname)

	fmt.Println("Url:", url)

	req, err := http.NewRequest("PUT", url, contents)
	if err != nil {

		return err
	}

	fmt.Println("Content-Type: ", mimeType)
	req.Header.Add("Content-Type", mimeType)

	resp, err := this.Client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {

		return errors.New(resp.Status)
	}

	return nil
}

func (this *Storage) Upload(bucket Bucket, objectname string, contents io.ReadSeeker) error {

	mimeType := DetectContentType(contents)
	if len(mimeType) == 0 {
		return errors.New("Fail to find mimeType")
	}

	return this.UploadMime(bucket, objectname, contents, mimeType)
}

func (this *Storage) DownloadBytes(bucket Bucket, objectname string) ([]byte, error) {

	url := fmt.Sprintf("%s/%s/%s/%s", this.Config.Host, this.Config.User, bucket, objectname)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {

		return nil, err
	}

	resp, err := this.Client.Do(req)
	if err != nil {

		return nil, err
	}

	if resp.StatusCode >= 300 {

		return nil, errors.New(resp.Status)
	}

	defer resp.Body.Close()

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return contents, nil

}
