package main

import (
	"fmt"
	"os"
	"testing"
	"text/template"
)

type Person struct {
	Name string
	Age  int
}

func TestTemplate(t *testing.T) {
	tmp, err := template.ParseFiles("./index.html")
	if err != nil {
		glog.Info("parse file err:", err)
		return
	}
	p := Person{Name: "Mary", Age: 11}
	if err := tmp.Execute(os.Stdout, p); err != nil {
		glog.Info("There was an error:", err.Error())
	}
}
