package main

import (
	"github.com/labstack/echo"
	"encoding/json"
	"io/ioutil"
	"path"
	"os"
	"io"
	"path/filepath"
	"log"
	"fmt"
	"net/http"
)

type Config struct {
	Server struct {
		Port        int `json:"port"`
	} `json:"server"`
	DB struct {
		Host        string `json:"host"`
		DB          string `json:"db"`
	} `json:"database"`
}
var config Config

var mongo MongoConnection

func copyFile(path, dest string) error {
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	e := out.Close()
	if e != nil {
		return e
	}

	return nil
}

var AppPath string
func init() {
	var err error
	AppPath, err = filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}

	if _, err = os.Stat(path.Join(AppPath, "config.json")); os.IsNotExist(err) {
		copyFile(path.Join(AppPath, "resources", "config.json"), path.Join(AppPath, "config.json"))
	}

	b, err := ioutil.ReadFile(path.Join(AppPath, "config.json"))
	if err != nil {
		log.Panic(err)
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	mongo = MongoConnection{}
	err := mongo.Init(config.DB.Host, config.DB.DB)
	if err != nil {
		panic(err)
	}

	e := echo.New()

	fs := http.FileServer(http.Dir(path.Join(AppPath, "kakin")))
	e.GET("/*", echo.WrapHandler(fs))

	e.POST("/api/v1/verify", VerifyAccount)
	e.POST("/api/v1/account", GetAccount)
	e.POST("/api/v1/renew", RenewToken)
	e.POST("/api/v1/usegold", SubtractGold)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Server.Port)))
}
