package Conf

import (
	"testing"
    "fmt"
)

func Test(t *testing.T) {
	cf := InitConf(&Load{
		Default:  false,
		SpecPath: "./",
		FileName: "example_conf.yaml",
	},
		NewYamlParser())
	fmt.Println(cf.Int("b"))
	fmt.Println(cf.String("a"))
	fmt.Println(cf.Ints("c"))
	fmt.Println(cf.Strings("d"))
	fmt.Println(cf.Strings("e"))
    fmt.Println(cf.Strings("string"))
}
