package main

import (
	"fmt"

	"github.com/baa-middleware/recovery"
	"github.com/baa-middleware/static"
	"github.com/go-baa/baa"
	"github.com/timest/env"
)

type config struct {
	AK        string
	SK        string
	Bucket    string
	IsImage   bool `default:"false"`
	BaseURI   string
	MaxSize   int64  `default:"2097152"`
	AllowType string `default:".jpg;.png;.jpeg;"`
}

var (
	accessKey     string
	secretKey     string
	bucket        string
	bucketisimage bool
	maxsize       int64
	allowtype     string
	baseuri       string
)

func main() {
	cfg := new(config)
	err := env.Fill(cfg)
	if err != nil {
		panic(err)
	}

	accessKey = cfg.AK
	secretKey = cfg.SK
	bucket = cfg.Bucket
	bucketisimage = cfg.IsImage
	maxsize = cfg.MaxSize
	allowtype = cfg.AllowType
	baseuri = cfg.BaseURI

	if accessKey == "" || secretKey == "" {
		fmt.Println("please setup accessKey and secretKey first")
		return
	}

	app := baa.New()
	app.Use(recovery.Recovery())
	app.Use(static.Static("/assets", "public/assets", false, nil))
	app.Use(static.Static("/html", "public/html", false, nil))

	app.Get("/", Index)
	app.Post("/upload", doUpload)
	app.Get("/refresh", refresh)
	app.Post("/refresh", refresh)
	app.Post("/delete", delete)
	app.Get("/getbandwidthdata", getBandwidthData)
	app.Get("/getfluxdata", getFluxData)

	app.Run(":80")

}
