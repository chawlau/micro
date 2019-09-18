package main

import (
	"html/template"
	"os"
	"path"

	"github.com/golang/glog"
)

type RouterGenerator struct {
}

func (d *RouterGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	filename := path.Join(opt.Output, "router/router.go")

	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		glog.Info("open file failed ", filename, " err ", err)
		return
	}

	defer file.Close()

	err = d.render(file, router_template, metaData)

	if err != nil {
		glog.Info("render failed err ", err)
		return
	}
	return
}

func (d *RouterGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main")

	t, err = t.Parse(data)
	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func init() {
	gen := &RouterGenerator{}
	Register("RouterGenerator", gen)
}
