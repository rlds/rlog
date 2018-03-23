//
//  
//  rlog
//
//  Created by 吴道睿 on 17/2/18.
//  Copyright © 2017年 吴道睿. All rights reserved.
//

package rlog

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)


// 可用于创建日志文件的文件夹
var logDirs []string

// 通过参数传递的设置值默认为 ../log/
var logDir = "./log/"
var logDirMonitor = "./log/monitor"

var (
	pid      = os.Getpid()               //pid
	program  = filepath.Base(os.Args[0]) //程序的名字
	removeLogFilePrevDay  = 9            //默认删除日志日期 默认 9天
)

// 生成文件名称 及文件快捷链接名
func logName(tag string, t time.Time) (name, link string) {
	name = fmt.Sprintf("%s.%s_%04d%02d%02d_%02d.%02d_%d.log",
		program,
		tag,
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		pid)
	return name, program + "." + tag +".log"
}

//创建文件 同时建立一个连接 成功返回文件 和文件名 忽略错误
func create(tag string, t time.Time) (f *os.File, filename string, err error) {
	if len(logDir) == 0 {
		return nil, "", errors.New("log: no log dirs")
	}
	name, link := logName(tag, t)
	var lastErr error
	fname := filepath.Join(logDir, name)
	f, err = os.Create(fname)
	if err == nil {
		symlink := filepath.Join(logDir, link)
		os.Remove(symlink)        // 忽略错误
		os.Symlink(name, symlink) // 忽略错误
		
		//清理日志文件
		go removeLogFiles(logDir)
		
		return f, fname, nil
	}
	lastErr = err
	return nil, "", fmt.Errorf("log: cannot create log: %v", lastErr)
}

//创建文件 同时建立一个连接 成功返回文件 和文件名 忽略错误
func create_monitor(tag string, t time.Time) (f *os.File, filename string, err error) {
	if len(logDirMonitor) == 0 {
		return nil, "", errors.New("log: no log dirs")
	}
	name, link := logName(tag, t)
	var lastErr error
	fname := filepath.Join(logDirMonitor, name)
	f, err = os.Create(fname)
	if err == nil {
		symlink := filepath.Join(logDirMonitor, link)
		os.Remove(symlink)        // 忽略错误
		os.Symlink(name, symlink) // 忽略错误
		
		//清理日志文件
		go removeLogFiles(logDirMonitor)
		
		return f, fname, nil
	}
	lastErr = err
	return nil, "", fmt.Errorf("log: cannot create log: %v", lastErr)
}
