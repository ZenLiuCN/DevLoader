package main

import (
	"os"
	"path/filepath"
	"strings"
	"path"
	"strconv"
	"time"

	"log"
	"io"
	"io/ioutil"
	"bytes"
	"github.com/axgle/mahonia"
	"runtime"
	"os/exec"
	"fmt"
	"syscall"
	"errors"
	"gopkg.in/yaml.v2"
	"encoding/json"
	"runtime/debug"
)

//-ldflags "-H windowsgui"

var logFile string
var logger *log.Logger
var exeRoot string

const (
	_s  = 1000000000
	_ms = 1000000
	_us = 1000
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			if logger == nil {
				f, e := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
				if e != nil {
					println(e.Error())
					panic(e.Error())
				}
				defer f.Close()
				initLogger(f)

			}
			switch r.(type) {
			case error:
				logger.Fatalf("system error \n %s \n%s", r.(error).Error(), string(debug.Stack()))
			case string:
				logger.Fatalf("system error \n %s\n%s", r, string(debug.Stack()))
			default:
				logger.Fatal("system error \n ", r, "\n", string(debug.Stack()))
			}
		}
	}()
	conf := getVars()
	if !fileExist(conf) {
		f, e := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
		if e != nil {
			panic(e.Error())
		}
		defer f.Close()
		initLogger(f)
		logger.Printf("will create config template:%s", conf)
		er := ioutil.WriteFile(conf, []byte(configTemplate), 0644)
		if er != nil {
			logger.Fatalf("create config template %s failed \n %s", conf, er.Error())
		}
		logger.Printf("create config template %s done!", conf)
		os.Exit(0)
	} else if config, e := ParseFile(conf); e != nil {
		fmt.Printf("error to parse config %s \n %s", conf, e.Error())
		f, e1 := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
		if e1 != nil {
			panic(e1.Error())
		}
		defer f.Close()
		initLogger(f)
		logger.Printf("parse config %s error :%s", conf, e.Error())
		os.Exit(1)
	} else {
		d, e := json.Marshal(config)
		fmt.Printf(" parseed config is: \n %s", string(d))
		if config.Debug {
			f, e := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
			if e != nil {
				panic(e.Error())
			}
			defer f.Close()
			defer func() { time.Sleep(1500 * _ms) }()
			initLogger(f)
		} else {
			initLogger(os.Stderr)
		}
		logger.Printf(" parseed config is: \n %s", string(d))
		er := config.Process()
		if er != nil {
			//fmt.Printf("error to process config %s \n %s", conf, e.Error())
			logger.Fatalf("process config error :%s", e.Error())
			os.Exit(1)
		}
		defer func() {
			logger.Printf(" will excute RunQuit")
			config.RunQuit.Exec()
		}()

	}

}

func getVars() string {
	exeFullPath, _ := filepath.Abs(os.Args[0])
	exeRoot = filepath.Dir(exeFullPath)
	_, exeName := filepath.Split(exeFullPath)
	ext := filepath.Ext(exeFullPath)
	name := strings.Replace(exeName, ext, "", -1)
	logFile = path.Join(exeRoot, name+"-"+strconv.FormatInt(time.Now().Unix(), 10)) + ".log"
	config := path.Join(exeRoot, name) + ".yml"
	if fileExist(path.Join(exeRoot, name) + ".yaml") {
		config = path.Join(exeRoot, name) + ".yaml"
	}
	//fmt.Printf("inti vars %s\n%s\n%s\n%s\n%s\n", exeFullPath, exeName, name, logFile, config)
	return config
}

func fileExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func initLogger(out io.Writer) {
	logger = log.New(out, "", log.LstdFlags|log.Lshortfile)
}

const (
	pathListSep    = string(os.PathListSeparator)
	configTemplate = `# 文件编码 utf-8
# debug             输出日志 默认 false 开启将输出到 程序名-时间.log 文件中,否则输出到StdErr
# runbefore         配置环境变量前运行
# env               要配置的环境变量数组
# runafter          配置环境变量后运行
# runquit           退出前运行
# *.x86             默认运行（x86 架构或者没有x64配置）
# *.x64             x64 架构运行内容,不配置 则运行x86
# run*.x*.path      命令 相对启动目录用 "/" 开始
# run*.x*.param     参数字符串数组
# run*.x*.show      是否显示界面 默认false
# run*.x*.must      是否必须运行成功 默认false
# run*.x*.fixparam  是否修复参数的相对目录 "./"开始 默认false
# run*.x*.wait      是否等待运行完毕 默认false
# run*.x*.sleep     等待时间 秒 ;默认 0 ,如不为 0 ,可以设置 path为空
# env.*.name        环境变量名
# env.*.type        配置类型 0 默认 追加路径 1 前添加路径 2 替换路径 3 追加值 4 前添加值 5 替换值 6 新建值 7 新建路径
debug: false
runbefore:
    x86:
        -   path: notepad
            param:
                - execute1
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
        -   path: notepad
            param:
                - execute2
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
    x64:
        -   path: notepad
            param:
                - execute64-1
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
        -   path: notepad
            param:
                - execute64-2
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
env:
    -   name: path
        type: 0
        x86:
            - /dir
        x64:
            - /dir
runafter:
    x86:
        -   path: notepad
            param:
                - abcd
            show: true
            must: true
            fixparam: true
            wait: false
            sleep: 2
`
)

type Config struct {
	Debug     bool    `yaml:",omitempty"`
	RunBefore Execute `yaml:",omitempty"`
	Env       []Env   `yaml:",omitempty"`
	RunAfter  Execute `yaml:",omitempty"`
	RunQuit   Execute `yaml:",omitempty"`
}

func (c *Config) Process() (er error) {
	defer func() {
		if r := recover(); r != nil {
			if logger == nil {
				f, e := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_SYNC, 0644)
				if e != nil {
					println(e.Error())
					panic(e.Error())
				}
				defer f.Close()
				initLogger(f)

			}
			switch r.(type) {
			case error:
				logger.Fatalf("process config Failed \n %s \n%s", r.(error).Error(), string(debug.Stack()))
			case string:
				logger.Fatalf("process config Failed \n %s\n%s", r, string(debug.Stack()))
			default:
				logger.Fatal("process config Failed \n ", r, "\n", string(debug.Stack()))
			}
		}
	}()
	if !c.RunBefore.IsZero() {
		logger.Printf("will execute RunBefore : %s", c.RunBefore.String())
		er = c.RunBefore.Exec()
		if er != nil {
			logger.Fatalf("execute RunBefore error: %s", er.Error())
			return
		}
	}
	if len(c.Env) != 0 {
		for _, v := range c.Env {
			er = v.Set()
			if er != nil {
				return
			}
		}
	}
	if !c.RunAfter.IsZero() {
		logger.Printf("will execute RunAfter : %s", c.RunAfter.String())
		er = c.RunAfter.Exec()
		if er != nil {
			return
		}
	}
	return
}
func (c *Config) String() string {
	d, e := json.Marshal(c)
	if e != nil {
		panic(e.Error())
	}
	return string(d)
}

/**
	type: 0 path append 1 path preappend 2 path replace 3 value append 4 value preappend 5 value replace 6 create if not Exists
 */
type Env struct {
	Name string
	Type int      `yaml:",omitempty"`
	X86  []string `yaml:",omitempty"`
	X64  []string `yaml:",omitempty"`
}

func (c *Env) String() string {
	d, e := json.Marshal(c)
	if e != nil {
		panic(e.Error())
	}
	return string(d)
}
func (e *Env) Set() error {
	if len(e.Name) == 0 {
		return errors.New("env name must not empty")
	}
	if strings.ContainsAny(e.Name, " \r\n\t") {
		return errors.New("env name must not have white chars")
	}
	if len(e.X64) > 0 && runtime.GOARCH == "amd64" {
		return setEnvs(e.Type, e.Name, e.X64)
	} else if len(e.X86) > 0 {
		return setEnvs(e.Type, e.Name, e.X86)
	}
	return nil
}

type Execute struct {
	X86 []Excuteable `yaml:",omitempty"`
	X64 []Excuteable `yaml:",omitempty"`
}

func (c *Execute) String() string {
	d, e := json.Marshal(c)
	if e != nil {
		panic(e.Error())
	}
	return string(d)
}
func (e *Execute) IsZero() bool {
	return len(e.X86) == 0 && len(e.X64) == 0
}
func (e *Execute) Exec() (er error) {
	if runtime.GOARCH == `amd64` && e.X64 != nil && len(e.X64) > 0 {
		for _, e := range e.X64 {
			er = e.Exec()
			if er != nil {
				return
			}
		}
	} else if e.X86 != nil && len(e.X86) > 0 {
		for _, e := range e.X86 {
			er = e.Exec()
			if er != nil {
				return
			}
		}
	}
	return
}

type Excuteable struct {
	Path     string   `yaml:",omitempty"`
	Param    []string `yaml:",omitempty"`
	Show     bool     `yaml:",omitempty"`
	Must     bool     `yaml:",omitempty"`
	FixParam bool     `yaml:",omitempty"`
	Sleep    int      `yaml:",omitempty"`
	Wait     bool     `yaml:",omitempty"`
}

func (c *Excuteable) String() string {
	d, e := json.Marshal(c)
	if e != nil {
		panic(e.Error())
	}
	return string(d)
}
func (c *Excuteable) Exec() (er error) {
	if len(c.Path) == 0 {
		if c.Sleep != 0 {
			logger.Printf("will sleep %ds", c.Sleep)
			time.Sleep(time.Duration(c.Sleep * 1000000000))
			logger.Printf("sleeped %ds", c.Sleep)
			return nil
		}
		if c.Must {
			logger.Fatal("Excuteable path must not empty")
		}

		return errors.New("Excuteable path must not empty")
	}

	if strings.HasPrefix(c.Path, "/") {
		c.Path = fixPath(c.Path)
	} else if runtime.GOOS == `windows` && !strings.Contains(c.Path, ":/") {
		c.Path, er = exec.LookPath(c.Path)
		if er != nil {
			er = errors.New(fmt.Sprintf("Can't locate excuteable of '%s' via %s.", c.String(), er.Error()))
			if c.Must {
				logger.Fatal(er.Error())
			}
			return
		}
	}
	if c.FixParam && c.Param != nil {
		for i, p := range c.Param {
			if strings.Contains(p, "./") {
				c.Param[i] = fixPath2(p)
			}
		}
	}
	if c.Param == nil {
		c.Param = []string{}
	}
	cmd := exec.Command(c.Path, c.Param...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: !c.Show}
	var errs bytes.Buffer
	cmd.Stderr = &errs
	logger.Printf("will execute: %s %s ", cmd.Path, strings.Join(cmd.Args, " "))
	if c.Wait {
		if err := cmd.Run(); err != nil {
			errs := errors.New(fmt.Sprintf("error execute: \n %s %s \n erro:%s \n output:%s \n%s", cmd.Path, strings.Join(cmd.Args, " "), err.Error(), fromGBK(errs), c.String()))
			if c.Must {
				logger.Fatal(errs.Error())
			}
			return
		}
		logger.Printf("execute: \n %s %s done! \n output:%s", cmd.Path, strings.Join(cmd.Args, " "), fromGBK(errs))
	} else {
//		go func() {
			if err := cmd.Start(); err != nil {
				errr := errors.New(fmt.Sprintf("error execute: \n %s %s \n erro:%s \n output:%s \n%s", cmd.Path, strings.Join(cmd.Args, " "), err.Error(), fromGBK(errs), c.String()))
				if c.Must {
					logger.Fatal(errr.Error())
				}
			}
			logger.Printf("execute: \n %s %s done! \n output:%s", cmd.Path, strings.Join(cmd.Args, " "), fromGBK(errs))
//		}()
	}
	if c.Sleep != 0 {
		logger.Printf("will sleep %ds", c.Sleep)
		time.Sleep(time.Duration(c.Sleep * _s))
		logger.Printf("sleeped %ds", c.Sleep)
		return nil
	}
	return nil
}


/**
parse configuation from file
 */
func ParseFile(file string) (c *Config, e error) {
	var f []byte
	f, e = ioutil.ReadFile(file)
	if e != nil {
		fmt.Printf("error to read file %s \n %s", file, e.Error())
		return nil, e
	}
	c = new(Config)
	e = yaml.Unmarshal(f, c)
	if e != nil {
		fmt.Printf("error to parse config %s \n %s", file, e.Error())
		return nil, e
	}
	return c, nil
}
/**
parse configuation from bytes
 */
func Parse(bytes []byte) (c *Config, e error) {
	c = new(Config)
	e = yaml.Unmarshal(bytes, c)
	if e != nil {
		return nil, e
	}
	return c, nil
}
/**
convert / path to abs path
 */
func fixPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return filepath.Join(exeRoot, strings.Replace(path, "/", "", 1))
	}
	return path
}
/**
convert ./path to abs path
 */
func fixPath2(path string) string {
	p, e := filepath.Abs(path)
	if e != nil {
		return path
	}
	return p
}
/**
setting env
 */
func appendEnvPath(name, path string) error {
	src:=os.Getenv(name)
	tar:=fixPath(path)
	if len(src)!=0{
		tar=src+pathListSep+tar
	}
	return os.Setenv(name, tar)
}
func preappendEnvPath(name, path string) error {
	src:=os.Getenv(name)
	tar:=fixPath(path)
	if len(src)!=0{
		tar=tar+pathListSep+src
	}
	return os.Setenv(name, tar)
}
func replaceEnvPath(name, path string) error {
	return os.Setenv(name, fixPath(path))
}
func createEnvPath(name, path string) (error, bool) {
	if len(os.Getenv(name)) == 0 {
		return os.Setenv(name, fixPath(path)), true
	}
	return nil, false

}
func appendEnvValue(name, value string) error {
	src:=os.Getenv(name)
	tar:=value
	if len(src)!=0{
		tar=src+pathListSep+tar
	}
	return os.Setenv(name, tar)
}
func preappendEnvValue(name, value string) error {
	src:=os.Getenv(name)
	tar:=value
	if len(src)!=0{
		tar=tar+pathListSep+src
	}
	return os.Setenv(name, tar)
}
func replaceEnvValue(name, value string) error {
	return os.Setenv(name, value)
}
func createEnvValue(name, value string) (error, bool) {
	if len(os.Getenv(name)) == 0 {
		return os.Setenv(name, value), true
	}
	return nil, false

}
/**
create env
 */
func setEnvs(mode int, name string, value []string) error {
	switch mode {
	case 1:
		for _, v := range value {
			if e := preappendEnvPath(name, v); e != nil {
				logger.Fatalf("set env %s-->%s failed \n  %s", name, v, e.Error())
				return e
			}
		}
	case 2:
		replaceEnvValue(name, "")
		for _, v := range value {
			if e := appendEnvPath(name, v); e != nil {
				logger.Fatalf("set env %s-->%s failed \n  %s", name, v, e.Error())
				return e
			}
		}
	case 3:
		for _, v := range value {
			if e := appendEnvValue(name, v); e != nil {
				logger.Fatalf("set env %s-->%s failed \n  %s", name, v, e.Error())
				return e
			}
		}
	case 4:
		for _, v := range value {
			if e := preappendEnvValue(name, v); e != nil {
				logger.Fatalf("set env %s-->%s failed \n  %s", name, v, e.Error())
				return e
			}
		}
	case 5:
		replaceEnvValue(name, "")
		for _, v := range value {
			if e := preappendEnvValue(name, v); e != nil {
				logger.Fatalf("set env %s-->%s failed \n  %s", name, v, e.Error())
				return e
			}
		}
	case 6:
		if e, s := createEnvValue(name, ""); !s {
			logger.Printf("skip env %s cause it's exists", name)
			return nil
		} else if e != nil {
			logger.Fatalf("create env %s failed \n %s", name, e.Error())
			return e
		}
		for _, v := range value {
			if e := preappendEnvValue(name, v); e != nil {
				logger.Fatalf("set env %s-->%s failed \n  %s", name, v, e.Error())
				return e
			}
		}
	case 7:
		if e, s := createEnvPath(name, ""); !s {
			logger.Printf("skip env %s cause it's exists", name)
			return nil
		} else if e != nil {
			logger.Fatalf("create env %s failed \n %s", name, e.Error())
			return e
		}
		for _, v := range value {
			if e := preappendEnvValue(name, v); e != nil {
				logger.Fatalf("set env %s-->%s failed \n  %s", name, v, e.Error())
				return e
			}
		}
	default:
		for _, v := range value {
			if e := appendEnvPath(name, v); e != nil {
				logger.Fatalf("set env %s-->%s failed \n  %s", name, v, e.Error())
				return e
			}
		}
	}
	return nil
}
/**
translate NONE Unicode system output
 */
func fromGBK(buffer bytes.Buffer) string {
	defer func() {
		if r := recover(); r != nil {
			switch r.(type) {
			case error:
				logger.Fatalf("convert GBK to UTF8 Failed \n %s \n%s", r.(error).Error(), string(debug.Stack()))
			case string:
				logger.Fatalf("convert GBK to UTF8 Failed \n %s\n%s", r, string(debug.Stack()))
			default:
				logger.Fatal("convert GBK to UTF8 Failed \n ", r, "\n", string(debug.Stack()))
			}
		}
	}()
	enc := mahonia.NewDecoder(`GBK`)
	return enc.ConvertString(buffer.String())
}
