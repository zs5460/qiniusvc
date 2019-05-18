// Package static provider a baa middleware for static file before router.
package static

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"strings"

	"gopkg.in/baa.v1"
)

// compatible with go net standard indexPage
const indexPage = "/index.html"

type static struct {
	handler baa.HandlerFunc
	prefix  string
	dir     string
	index   bool
}

// Static returns a baa middleware with static file before router.
func Static(prefix, dir string, index bool, h baa.HandlerFunc) baa.HandlerFunc {
	if len(prefix) > 1 && prefix[len(prefix)-1] == '/' {
		prefix = prefix[:len(prefix)-1]
	}
	if len(dir) > 1 && dir[len(dir)-1] == '/' {
		dir = dir[:len(dir)-1]
	}
	s := &static{
		dir:     dir,
		index:   index,
		prefix:  prefix,
		handler: h,
	}

	return func(c *baa.Context) {
		file := c.Req.URL.Path
		if strings.HasPrefix(file, s.prefix) {
			file = file[len(s.prefix):]
		} else {
			c.Next()
			return
		}

		if len(file) > 0 && file[0] == '/' {
			file = file[1:]
		}
		file = s.dir + "/" + file

		if s.handler != nil {
			s.handler(c)
		}

		// directory index check
		if f, err := os.Stat(file); err == nil {
			if f.IsDir() {
				if s.index {
					// if no end slash, add slah and redriect
					if c.Req.URL.Path[len(c.Req.URL.Path)-1] != '/' {
						c.Redirect(302, c.Req.URL.Path+"/")
						c.Break()
						return
					}
					listDir(file, s, c)
				} else {
					// check index
					if err := serveIndex(file+indexPage, c); err != nil {
						c.Resp.WriteHeader(http.StatusForbidden)
					}
				}
				c.Break()
				return
			}
			// file
			if strings.HasSuffix(file, indexPage) {
				if err := serveIndex(file, c); err != nil {
					c.Error(err)
				}
			} else {
				http.ServeFile(c.Resp, c.Req, file)
			}
			c.Break()
		}

		c.Next()
	}
}

// listDir list given dir files
func listDir(dir string, s *static, c *baa.Context) {
	f, err := os.Open(dir)
	if err != nil {
		c.Error(fmt.Errorf("baa.Static listDir Error: %s", err))
	}
	defer f.Close()
	fl, err := f.Readdir(-1)
	if err != nil {
		c.Error(fmt.Errorf("baa.Static listDir Error: %s", err))
	}

	dirName := f.Name()
	dirName = dirName[len(s.dir):]
	c.Resp.Header().Add("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(c.Resp, "<h3 style=\"padding-bottom:5px;border-bottom:1px solid #ccc;\">%s</h3>\n", dirName)
	fmt.Fprintf(c.Resp, "<pre>\n")
	var color, name string
	for _, v := range fl {
		name = v.Name()
		color = "#333333"
		if v.IsDir() {
			name += "/"
			color = "#3F89C8"
		}
		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		url := url.URL{Path: name}
		fmt.Fprintf(c.Resp, "<a style=\"color:%s\" href=\"%s\">%s</a>\n", color, url.String(), template.HTMLEscapeString(name))
	}
	fmt.Fprintf(c.Resp, "</pre>\n")
}

func serveIndex(file string, c *baa.Context) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()
	fs, err := f.Stat()
	if err != nil {
		return err
	}
	http.ServeContent(c.Resp, c.Req, f.Name(), fs.ModTime(), f)
	return nil
}
