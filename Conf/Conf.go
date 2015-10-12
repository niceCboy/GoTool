package Conf

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type (
	Conf struct {
		m map[string]interface{}
		l *Load
		p Parser
	}

	Load struct {
		Default  bool //默认标记，真时RunMode值有效，否则SpecPath值有效
		RunMode  string
		SpecPath string
		FileName string
	}

	Parser interface {
		Unmarshal(data []byte, v interface{}) error
	}
)

func InitConf(l *Load, p Parser) *Conf {
	cf := &Conf{
		m: make(map[string]interface{}),
		l: l,
		p: p,
	}

	//Loading config file...
	var (
		stream []byte
		e      error
	)
	if l.Default {
		stream, e = l.loadDefault()
	} else {
		stream, e = l.loadSpec()
	}
	if e != nil {
		panic(e)
	}

	//Parse byte stream
	if parseError := p.Unmarshal(stream, &cf.m); parseError != nil {
		panic(parseError)
	}

	return cf
}

//默认情况下，从执行文件的当前路径开始查找 conf/(runmode)/(filename)的文件
func (l *Load) loadDefault() (ret []byte, e error) {
	if l.RunMode == "" || l.FileName == "" {
		e = errors.New("lack of Load config..")
		return
	}
	path, _ := os.Getwd()
loadloop:
	for {
		if _, err := os.Stat(path + "/conf/" + l.RunMode + "/" + l.FileName); err == nil {
			ret, e = ioutil.ReadFile(l.SpecPath)
			return
		}
		path = path[:strings.LastIndex(path, "/")] //上一级目录路径
		if path == "/" {                           //已遍历到根目录
			e = errors.New("can't find the config file.")
			break loadloop
		}
	}
}

func (l *Load) loadSpec() ([]byte, error) {
	return ioutil.ReadFile(l.SpecPath)
}
