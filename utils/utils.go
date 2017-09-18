package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os/exec"
	"reflect"
	"regexp"
	"strings"
	"time"

	"errors"

	"fmt"

	"github.com/tidwall/gjson"
)

var configJson string = ""

//每整点小时调用
func SetTimerPerHour(F func()) {
	go func() {
		for {
			F()
			now := time.Now()
			next := now.Add(time.Hour)
			next = time.Date(next.Year(), next.Month(), next.Day(), next.Hour(), 0, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(now))
			<-t.C
		}
	}()
}

//Md5校验
func Md5Check(checkStr string, gen string) bool {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(gen))
	cipherStr := md5Ctx.Sum(nil)
	token := hex.EncodeToString(cipherStr)
	if token != checkStr {
		return false
	}
	return true
}

func CreateMd5(gen string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte("cgyx2017"))
	cipherStr := md5Ctx.Sum(nil)
	return hex.EncodeToString(cipherStr)
}

func ExeShell(syscmd string, dir string, args string) (string, error) {
	log.Println("begin execute shell.....", syscmd, dir, args)
	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	cmd := exec.Command(syscmd, dir, args) //"GameConfig/gitCommit", "zoneo")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("stdoutpipe:", err.Error())
		return "", err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Println("cmd start:", err.Error())
		return "", err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ioutil.ReadAll:", err.Error())
		return "", err
	}
	e := cmd.Wait()
	if e != nil {
		log.Println("Exeshell error:", e.Error())
	}
	s := strings.Replace(string(opBytes), "\n", "", -1)
	log.Println(string(opBytes))
	return s, nil
}

func ExeShellArgs2(syscmd string, dir string, arg1 string, arg2 string) (string, error) {
	log.Println("begin execute shell.....", syscmd, dir, "arg1:", arg1, "arg2:", arg2)
	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	cmd := exec.Command(syscmd, dir, arg1, arg2) //"GameConfig/gitCommit", "zoneo")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("stdoutpipe:", err.Error())
		return "", err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Println("cmd start:", err.Error())
		return "", err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ioutil.ReadAll:", err.Error())
		return "", err
	}
	e := cmd.Wait()
	if e != nil {
		log.Println("Exeshell error:", e.Error())
	}
	s := strings.Replace(string(opBytes), "\n", "", -1)
	return s, nil
}

func ExeShellArgs3(syscmd string, dir string, arg1 string, arg2 string, arg3 string) (string, error) {
	log.Println("begin execute shell.....", syscmd, dir, "arg1:", arg1, "arg2:", arg2, "arg3:", arg3)
	// 执行系统命令
	// 第一个参数是命令名称
	// 后面参数可以有多个，命令参数
	cmd := exec.Command(syscmd, dir, arg1, arg2, arg3) //"GameConfig/gitCommit", "zoneo")
	// 获取输出对象，可以从该对象中读取输出结果
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Println("stdoutpipe:", err.Error())
		return "", err
	}
	// 保证关闭输出流
	defer stdout.Close()
	// 运行命令
	if err := cmd.Start(); err != nil {
		log.Println("cmd start:", err.Error())
		return "", err
	}
	// 读取输出结果
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("ioutil.ReadAll:", err.Error())
		return "", err
	}
	e := cmd.Wait()
	if e != nil {
		log.Println("Exeshell error:", e.Error())
	}
	s := strings.Replace(string(opBytes), "\n", "", -1)
	return s, nil
}

func MatchType(a string, f string) bool {
	reg := regexp.MustCompile(f)
	s := reg.FindAllString(a, -1)
	if len(s) > 0 {
		return true
	}
	return false
}

func AgentServiceType(agentName string, tmap map[int]string) int {
	for k, v := range tmap {
		if MatchType(agentName, v) {
			return k
		}
	}
	return 0
}

func LoadConfigJson() error {
	fBytes, err := ioutil.ReadFile("./config.json")
	if err != nil {
		return err
	}
	configJson = string(fBytes)
	return nil
}
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func GetConfigValue(key string, value interface{}) error {
	if configJson == "" {
		err := LoadConfigJson()
		if err != nil {
			return err
		}
	}
	v := gjson.Get(configJson, key)
	switch v.Type {
	case gjson.Null:
		return errors.New("none value")
	case gjson.String:
		DeepCopy(value, v.String())
	case gjson.Number:
		DeepCopy(value, int(v.Int()))
	case gjson.False:
	case gjson.True:
		DeepCopy(value, v.Bool())
	default:
		return errors.New("no result")
	}
	return nil
}

func GetConfigMap(key string, vMap interface{}) error {
	if configJson == "" {
		err := LoadConfigJson()
		if err != nil {
			return err
		}
	}

	result := gjson.Get(configJson, key)
	mapType := reflect.TypeOf(vMap)
	switch mapType.String() {
	case "*map[string]string":
		{
			s := make(map[string]string)
			for k, name := range result.Map() {
				s[k] = name.String()
			}
			DeepCopy(vMap, s)
		}
	case "*map[string]bool":
		{
			s := make(map[string]bool)
			for k, name := range result.Map() {
				s[k] = name.Bool()
			}
			DeepCopy(vMap, s)
		}
	case "*map[string]int":
		{
			s := make(map[string]int)
			for k, name := range result.Map() {
				s[k] = int(name.Int())
			}
			DeepCopy(vMap, s)
		}
	default:
		return errors.New(fmt.Sprintf("map type valid, %s", mapType.String()))
	}
	return nil
}

func GetConfigArray(key string, vArray interface{}) error {
	if configJson == "" {
		err := LoadConfigJson()
		if err != nil {
			return err
		}
	}
	result := gjson.Get(configJson, key)
	mapType := reflect.TypeOf(vArray)
	switch mapType.String() {
	case "*[]string":
		{
			s := []string{}
			for _, v := range result.Array() {
				s = append(s, v.String())
			}
			DeepCopy(vArray, s)
		}
	case "*[]bool":
		{
			s := []bool{}
			for _, v := range result.Array() {
				s = append(s, v.Bool())
			}
			DeepCopy(vArray, s)
		}
	case "*[]int":
		{
			s := []int{}
			for _, v := range result.Array() {
				s = append(s, int(v.Int()))
			}
			DeepCopy(vArray, s)
		}
	default:
		return errors.New("array type valid")
	}
	return nil
}
