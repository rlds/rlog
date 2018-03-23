//
//  loginit.go
//  rlog
//
//  Created by 吴道睿 on 17/2/18.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package rlog

import(
    "time"
)

/*
   日志输出到文件的初始化
*/

// 日志文件的最大大小.  最大 1800M
var MaxSize uint64 = 1024 * 1024 * 1800
var MaxBufferSize int = 1024 *  100
//可以主动调用的初始化
func LogInit(verbosity int32,logdir string,maxlogsize uint64,logbufferSize  int){
	if logbufferSize >0 {
		bufferSize = logbufferSize
	}else{
		bufferSize = MaxBufferSize
	}
	if maxlogsize > 0 {
		MaxSize    = maxlogsize
	}else{
		MaxSize    = 1024 * 1024 * 1800
	}
	if logdir != "" {
		logDir            = logdir
		
	}else{
		logDir     = "./"
	}
	//输出info日志级别的控制
	if verbosity > 0 {
		logging.verbosity.set(Level(verbosity))
		objLoging.verbosity.set(Level(verbosity))
	}else{
		logging.verbosity.set(4)
		objLoging.verbosity.set(4)
	}
	
	logging.stderrThreshold = infoLog
	objLoging.stderrThreshold = monitor
	//输出info日志级别的控制	
	go logging.flushDaemon()
}

//设置并启用monitor日志
func StartMonitorLog(logdir string){
	logDirMonitor     = logdir
	go objLoging.flushDaemon()
}

//定时把缓冲数据写入文件的 时间间隔
const flushInterval = 30 * time.Second

func (l *loggingT) flushDaemon() {
	for _ = range time.NewTicker(flushInterval).C {
		l.lockAndFlushAll()
	}
}

