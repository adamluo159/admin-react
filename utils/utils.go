package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
)

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
	md5Ctx.Write([]byte(gen))
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

func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}
