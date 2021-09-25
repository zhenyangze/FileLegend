package main

import (
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"
	//"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/text/gregex"
	//"github.com/gogf/gf/os/gcfg"
	"github.com/gogf/gf/os/gcmd"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/glog"
	"github.com/gogf/gf/os/gtime"
	//"github.com/gogf/gf/os/gtimer"
)

var wg sync.WaitGroup
var rootDir string
var rule string

var showlog string
var isforce string

func main() {
	lockFile := interface2String(g.Cfg().Get("pidpath"))
	lock, err := os.Create(lockFile)
	if err != nil {
		glog.Error("创建文件锁失败", err)
	}
	defer os.Remove(lockFile)
	defer lock.Close()

	err = syscall.Flock(int(lock.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		glog.Error("上一个任务未执行完成，暂停执行")
		os.Exit(0)
	}
	defer syscall.Flock(int(lock.Fd()), syscall.LOCK_UN)

	gcmd.BindHandle("test", testReg)
	gcmd.BindHandle("run", master)
	gcmd.AutoRun()
}

func checkLockFile() {
}

func master() {
	nodes := g.Cfg().Get("items.node")
	wg = sync.WaitGroup{}
	isforce = interface2String(g.Cfg().Get("isforce"))
	showlog = interface2String(g.Cfg().Get("showlog"))
	for k, _ := range nodes.([]interface{}) {
		index := strconv.Itoa(k)
		rootDir := interface2String(g.Cfg().Get("items.node." + index + ".rootdir"))
		rule := interface2String(g.Cfg().Get("items.node." + index + ".reg"))
		wg.Add(1)
		if !gfile.IsDir(rootDir) {
			glog.Error(rootDir + " is not dir")
			continue
		}
		go scanDir(rootDir, rule)
	}

	wg.Wait()

}
func testReg() {
	filePath := gcmd.GetOpt("file")
	rule := gcmd.GetOpt("rule")
	timeOut, err := getTimeOut(filePath, rule)
	if err != nil {
		return
	}
	glog.Println(timeOut)
}

func interface2String(inter interface{}) string {
	switch inter.(type) {
	case string:
		return inter.(string)
	case int:
		return strconv.Itoa(inter.(int))
	}
	return ""
}

// 遍历文件
func scanDir(dirPath, rule string) {
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			// 添加到任务中
			addTask(path, rule)
		}

		return nil
	})
	wg.Done()
}

func getTimeOut(filePath, rule string) (string, error) {
    fileContents := gfile.GetContents(filePath)
	match, _ := gregex.MatchString(rule, fileContents)

	if len(match) >= 2 {
		return match[1], nil
	}
	return "", nil
}

// 添加定时任务
func addTask(filePath, rule string) {
	// 根据规则获取过期时间
	timeOut, err := getTimeOut(filePath, rule)
	if err != nil {
		return
	}

	if showlog == "1" {
		glog.Println(filePath, timeOut)
	}

	expireTime, _ := strconv.ParseInt(timeOut, 0, 64)
	timeidff := expireTime - gtime.Now().Second()
	if (timeidff <= 0 && expireTime > 0) || (timeidff <= 0 && expireTime == 0 && isforce == "1") {
		// 直接删除
		runTask(filePath)
		return
	}
	// 获取时间差
	// 添加定时任务, 需要以任务名称为例
	//glog.Println(timeidff)
	//gtimer.SetTimeout(time.Duration(timeidff)*1000*time.Millisecond, func() {
	//addTask(filePath, rule)
	//})
}

func delTask() {
	// 如果文件发生变动，或者时间发生变动，需要删除任务重新
}

// 到期处理
func runTask(filePath string) {
	if showlog == "1" {
		glog.Println("删除文件: " + filePath)
	}
	// 删除文件
	gfile.Remove(filePath)
}
