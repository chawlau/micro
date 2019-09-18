package main

import (
	"fmt"
	"html/template"
	"os"
	"path"

	"github.com/golang/glog"
	"github.com/ibinarytree/proto"
)

type CtrlGenerator struct{}

type RpcMeta struct {
	Rpc     *proto.RPC
	Package *proto.Package
	Prefix  string
}

func (d *CtrlGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	reader, err := os.Open(opt.Proto3Filename)

	if err != nil {
		glog.Info("open file ", opt.Proto3Filename, " failed ", err)
		return
	}

	defer reader.Close()
	return d.generateRpc(opt, metaData)
}

func (d *CtrlGenerator) generateRpc(opt *Option, metaData *ServiceMetaData) (err error) {
	for _, rpc := range metaData.Rpc {
		var file *os.File
		filename := path.Join(opt.Output, "controller", fmt.Sprintf("%s.go", rpc.Name))
		glog.Info("filename is ", filename)
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)

		if err != nil {
			glog.Info("open file failed", filename, " err ", err)
			return
		}

		rpcMeta := &RpcMeta{
			Package: metaData.Package,
			Rpc:     rpc,
			Prefix:  metaData.Prefix,
		}

		err = d.render(file, controller_template, rpcMeta)
		if err != nil {
			glog.Info("render controller failed err ", err)
			return
		}
		defer file.Close()
	}

	return
}

func (d *CtrlGenerator) render(file *os.File, data string, metaData *RpcMeta) (err error) {
	t := template.New("main")
	t, err = t.Parse(data)
	if err != nil {
		return
	}

	err = t.Execute(file, metaData)
	return
}

func init() {
	ctrl := &CtrlGenerator{}
	Register("ctrl generator", ctrl)
}
