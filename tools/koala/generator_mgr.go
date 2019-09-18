package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/glog"
	"github.com/ibinarytree/proto"
	"github.com/luci/go-render/render"
)

var genMgr *GeneratorMgr = &GeneratorMgr{
	genMap:   make(map[string]Generator),
	metaData: &ServiceMetaData{},
}

var AllDirList []string = []string{
	"controller",
	"idl",
	"main",
	"scripts",
	"conf/product",
	"conf/test",
	"app/router",
	"app/config",
	"model",
	"generate",
	"router",
}

type GeneratorMgr struct {
	genMap   map[string]Generator
	metaData *ServiceMetaData
}

func (g *GeneratorMgr) parseService(opt *Option) (err error) {
	reader, err := os.Open(opt.Proto3Filename)

	if err != nil {
		glog.Info("openfile failed ", opt.Proto3Filename, " err ", err)
		return
	}

	defer reader.Close()

	parser := proto.NewParser(reader)
	definition, err := parser.Parse()

	if err != nil {
		glog.Info("proto Parse failed ", opt.Proto3Filename, " err ", err)
		return
	}

	proto.Walk(definition,
		proto.WithService(g.handleService),
		proto.WithMessage(g.handleMessage),
		proto.WithRPC(g.handleRPC),
		proto.WithPackage(g.handlePackage),
	)

	return
}

func (g *GeneratorMgr) handleService(s *proto.Service) {
	g.metaData.Service = s
}

func (g *GeneratorMgr) handleMessage(m *proto.Message) {
	g.metaData.Messages = append(g.metaData.Messages, m)
}

func (g *GeneratorMgr) handleRPC(r *proto.RPC) {
	g.metaData.Rpc = append(g.metaData.Rpc, r)
}

func (g *GeneratorMgr) handlePackage(p *proto.Package) {
	g.metaData.Package = p
}

func (g *GeneratorMgr) createAllDir(opt *Option) (err error) {

	for _, dir := range AllDirList {
		fullDir := path.Join(opt.Output, dir)
		err = os.MkdirAll(fullDir, 0755)
		if err != nil {
			glog.Info("mkdir dir failed ", dir, " err ", err)
			return
		}
	}
	return
}

func (g *GeneratorMgr) initOutputDir(opt *Option) (err error) {
	goPath := ("/home/liuchao/Documents/GoCode")

	if len(opt.Prefix) > 0 {
		opt.Output = path.Join(goPath, "src", opt.Prefix)
		return
	}

	exeFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}

	lastIdx := strings.LastIndex(exeFilePath, "/")
	if lastIdx < 0 {
		err = fmt.Errorf("invalid path %v", exeFilePath)
		return
	}

	//opt.Output = strings.ToLower(exeFilePath[0:lastIdx])
	opt.Output = exeFilePath[0:lastIdx]
	srcPath := path.Join(goPath, "src/")

	if srcPath[len(srcPath)-1] != '/' {
		srcPath = fmt.Sprintf("%s/", srcPath)
	}

	opt.Prefix = strings.Replace(opt.Output, srcPath, "", -1)
	glog.Info("opt ", render.Render(opt), " goPath ", goPath)
	return
}

func (g *GeneratorMgr) Run(opt *Option) (err error) {
	err = g.initOutputDir(opt)
	if err != nil {
		return
	}

	err = g.parseService(opt)
	if err != nil {
		return
	}

	err = g.createAllDir(opt)
	if err != nil {
		return
	}

	g.metaData.Prefix = opt.Prefix
	for _, gen := range g.genMap {
		err = gen.Run(opt, g.metaData)
		if err != nil {
			return
		}
	}
	return
}

func Register(name string, gen Generator) (err error) {
	_, ok := genMgr.genMap[name]

	if ok {
		err = fmt.Errorf("generator %s is exists", name)
		return
	}

	genMgr.genMap[name] = gen
	return
}
