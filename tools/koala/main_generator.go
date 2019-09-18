package main

import (
	"html/template"
	"os"
	"path"

	"github.com/golang/glog"
)

type MainGenerator struct {
}

func (d *MainGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	fileName := path.Join(opt.Output, "main/main.go")

	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

	if err != nil {
		glog.Info("open file failed", fileName, " err ", err)
		return
	}

	defer file.Close()

	err = d.render(file, main_template, metaData)

	if err != nil {
		glog.Info(" render failed ", err)
		return
	}
	return
}

func (d *MainGenerator) render(file *os.File, data string, metaData *ServiceMetaData) (err error) {
	t := template.New("main")

	t, err = t.Parse(data)

	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func init() {
	dir := &MainGenerator{}
	Register("main generator", dir)
}
