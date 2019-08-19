//
//  remove_log_file.go
//  rlog
//
//  Created by 吴道睿 on 2017/10/16.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package rlog

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

/*
   删除n天前的日志文件
   删除规则：删除最近日期 前推n天之前的日志
*/

func GetRemoveLogFileDayNum() int {
	return removeLogFilePrevDay
}

func SetRemoveLogFileDayNum(n int) int {
	removeLogFilePrevDay = n
	return removeLogFilePrevDay
}

type S_Files []os.FileInfo

func (s S_Files) Len() int           { return len(s) }
func (s S_Files) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s S_Files) Less(i, j int) bool { return s[i].ModTime().Unix() > s[j].ModTime().Unix() }

func GetDirFiles(dir string) (files S_Files, err error) {
	files, err = ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	sort.Sort(files)
	return
}

// 获得需要删除的文件名列表
func GetRemoveFiles(dirpath, removeFilePrestr string, s S_Files, daynum int) (fpath []string) {
	if len(s) > 0 {
		f1 := s[0].ModTime()
		t := time.Date(f1.Year(), f1.Month(), f1.Day(), 0, 0, 0, 0, f1.Location())
		t = t.AddDate(0, 0, -daynum)
		utm := t.Unix()
		for _, f := range s {
			date := f.ModTime().Format("20060102")
			fname := f.Name()
			if strings.Contains(fname, removeFilePrestr) && strings.Contains(fname, date) { //仅清理包含日志日期标志的文件
				if f.ModTime().Unix() < utm {
					fpath = append(fpath, dirpath+"/"+fname)
				}
			}
		}
	}
	return
}

func removeLogFiles(dir string) (lasterr error) {
	//主动设置 <=0 则表示不删除日志 手动清理
	if removeLogFilePrevDay <= 0 {
		return
	}

	files, err := GetDirFiles(dir)
	if err != nil {
		lasterr = err
		fmt.Println("getDirFileErr:", err.Error())
		return
	}

	removeFilePreStr := program + "."

	rfiles := GetRemoveFiles(dir, removeFilePreStr, files, removeLogFilePrevDay)

	for _, file := range rfiles {
		err := os.Remove(file)
		if err != nil {
			lasterr = err
			//删除日志文件错误
			fmt.Println(file, "removeLogErr:", err.Error())
		}
	}
	return
}
