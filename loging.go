//
//  loging.go
//  rlog
//
//  Created by 吴道睿 on 17/2/18.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package rlog

import(
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

// loggingT 结构体
type loggingT struct {
	stderrThreshold severity // The -stderrthreshold 参数.
	freeList *buffer
	freeListMu sync.Mutex
	mu sync.Mutex
	file [numSeverity]flushSyncWriter
	verbosity Level      // -v 参数 日志输出级别
}


var logging loggingT

func (l *loggingT) getBuffer() *buffer {
	l.freeListMu.Lock()
	b := l.freeList
	if b != nil {
		l.freeList = b.next
	}
	l.freeListMu.Unlock()
	if b == nil {
		b = new(buffer)
	} else {
		b.next = nil
		b.Reset()
	}
	return b
}

func (l *loggingT) putBuffer(b *buffer) {
	if b.Len() >= 256 {
		// 大缓冲不处理
		return
	}
	l.freeListMu.Lock()
	b.next = l.freeList
	l.freeList = b
	l.freeListMu.Unlock()
}

//
func (l *loggingT) output(s severity, buf *buffer) {
	l.mu.Lock()
	data := buf.Bytes()
	if l.file[s] == nil {
		if err := l.createFiles(s); err != nil {
			os.Stderr.Write(data) // 确保消息会输出.
			l.exit(err)
		}
	}
	l.file[infoLog].Write(data)
	l.putBuffer(buf)
	l.mu.Unlock()
	if stats := severityStats[s]; stats != nil {
		atomic.AddInt64(&stats.lines, 1)
		atomic.AddInt64(&stats.bytes, int64(len(data)))
	}
}

//默认指定输出到标准输出
var logExitFunc func(error) = exitfunc 

func exitfunc(err error) {
	fmt.Println(err)
}

func (l *loggingT) exit(err error) {
	fmt.Fprintf(os.Stderr, "log: exiting because of error: %s\n", err)
	// 如果有定义logExitFunc 则执行已定义的
	if logExitFunc != nil {
		logExitFunc(err)
		return
	}
	l.flushAll()
	os.Exit(2)
}

//把内容存入磁盘
func (l *loggingT) lockAndFlushAll() {
	l.mu.Lock()
	l.flushAll()
	l.mu.Unlock()
}

func (l *loggingT) flushAll() {
	for s := monitor; s >= infoLog; s-- {
		file := l.file[s]
		if file != nil {
			file.Flush() // 忽略错误
			file.Sync()  // 忽略错误
		}
	}
}

