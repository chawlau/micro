package protocol

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/luci/go-render/render"
	"github.com/vmihailenco/msgpack"
)

type Person struct {
	Age    int     `json:"age"`
	Id     int64   `json:"id,string"`
	Name   string  `json:"name_xx,omitempty"`
	Salary float32 `json:"-"`
}

func TestMsgPack(t *testing.T) {
	var p = &Person{
		Age:    20,
		Id:     38888232322323222,
		Name:   "axx",
		Salary: 38822.2,
	}

	data, err := msgpack.Marshal(p)
	if err != nil {
		fmt.Printf("marshal failed, err:%v\n", err)
		return
	}

	ioutil.WriteFile("./msg.txt", data, 0777)

	data2, err := ioutil.ReadFile("./msg.txt")
	if err != nil {
		fmt.Printf("read file failed, err:%v\n", err)
		return

	}

	var person2 Person
	msgpack.Unmarshal(data2, &person2)
	fmt.Println("person2: ", render.Render(person2))
}
