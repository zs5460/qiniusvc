package main

import (
	"context"
	"errors"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/cdn"
	"github.com/qiniu/api.v7/v7/storage"
	"gopkg.in/baa.v1"
)

// Result ...
type Result struct {
	Error int    `json:"error"`
	URL   string `json:"url,omitempty"`
	Msg   string `json:"msg,omitempty"`
}

//Index ...
func Index(c *baa.Context) {
	c.Redirect(301, "/html")

}

func delete(c *baa.Context) {
	key := c.QueryTrim("key")
	if strings.HasPrefix(key, "http") {
		key = strings.Join(strings.Split(key, "/")[3:], "/")
	}
	mac := qbox.NewMac(accessKey, secretKey)
	cfg := storage.Config{
		// 是否使用https域名进行资源管理
		UseHTTPS: false,
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	err := bucketManager.Delete(bucket, key)
	if err != nil {
		c.Error(err)
		return
	}

	r := Result{Error: 0, Msg: "delete success"}
	c.JSON(200, r)

}

func doUpload(c *baa.Context) {
	file, header, err := c.GetFile("imgFile")
	if err != nil {
		c.Error(errors.New("no files to upload"))
		return
	}
	defer file.Close()

	if header.Size > maxsize {
		c.Error(errors.New("file is too large"))
		return
	}

	filename := strings.ToLower(header.Filename)
	extname := filepath.Ext(filename)
	if !strings.Contains(allowtype, extname+";") {
		c.Error(errors.New("file type is not allowed to upload"))
		return
	}

	newfilename := getNewFilename(extname)

	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(accessKey, secretKey)
	upToken := putPolicy.UploadToken(mac)
	formUploader := storage.NewFormUploader(nil)
	ret := storage.PutRet{}
	err = formUploader.Put(context.Background(),
		&ret,
		upToken,
		newfilename,
		file,
		header.Size,
		nil)
	if err != nil {
		c.Error(err)
		return
	}

	r := Result{Error: 0, Msg: "success", URL: baseuri + ret.Key}
	c.JSON(200, r)
}

func refresh(c *baa.Context) {
	urltorefresh := c.QueryTrim("resurl")
	if urltorefresh == "" {
		err := c.Redirect(307, "/html/refresh.html")
		if err != nil {
			c.Error(err)
		}
		return
	}

	if !strings.HasPrefix(urltorefresh, "http://") && !strings.HasPrefix(urltorefresh, "https://") {
		c.Error(errors.New("Invalid URL"))
		return
	}

	mac := qbox.NewMac(accessKey, secretKey)

	if bucketisimage {
		bucketManager := storage.NewBucketManager(mac, nil)
		key := strings.Join(strings.Split(urltorefresh, "/")[3:], "/")
		err := bucketManager.Delete(bucket, key)
		if err != nil {
			c.Error(err)
			return
		}
	}

	cdnManager := cdn.NewCdnManager(mac)
	urls := []string{urltorefresh}

	ret, err := cdnManager.RefreshUrls(urls)
	if err != nil {
		c.Error(err)
		return
	}

	r := Result{
		Error: 0,
		Msg:   "refresh " + urltorefresh + " " + ret.Error,
	}
	c.JSON(200, r)

}

//
func getBandwidthData(c *baa.Context) {
	mac := qbox.NewMac(accessKey, secretKey)
	cdnManager := cdn.NewCdnManager(mac)
	// domains := []string{
	// 	"images.rednet.cn",
	// 	"img.redimg.cn",
	// 	"img.app.rednet.cn",
	// 	"wmv2.rednet.cn"}
	domains := []string{}
	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")
	g := "5min"
	data, err := cdnManager.GetBandwidthData(startDate, endDate, g, domains)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, data)

}

func getFluxData(c *baa.Context) {
	mac := qbox.NewMac(accessKey, secretKey)
	cdnManager := cdn.NewCdnManager(mac)
	domains := []string{}
	startDate := time.Now().AddDate(0, -1, 0).Format("2006-01-02")
	endDate := time.Now().Format("2006-01-02")
	g := "day"
	data, err := cdnManager.GetFluxData(startDate, endDate, g, domains)
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(200, data)
}

func getNewFilename(extname string) string {
	rand.Seed(time.Now().UnixNano())
	return time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(899999)+100000) + extname

}
