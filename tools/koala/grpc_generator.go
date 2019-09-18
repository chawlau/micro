package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/golang/glog"
)

type GrpcGenerator struct{}

func (d *GrpcGenerator) Run(opt *Option, metaData *ServiceMetaData) (err error) {
	dir := path.Join(opt.Output, "generate", metaData.Package.Name)
	os.MkdirAll(dir, 0755)

	outputParams := fmt.Sprintf("plugins=grpc:%s/generate/%s", opt.Output, metaData.Package.Name)

	cmd := exec.Command("protoc", "--go_out", outputParams, opt.Proto3Filename)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err = cmd.Run()

	if err != nil {
		glog.Info("grpc generator failed err ", err)
		return
	}
	return
}

func init() {
	gc := &GrpcGenerator{}

	Register("grpc generator ", gc)
}
