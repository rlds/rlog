//
//  
//  rlog
//
//  Created by 吴道睿 on 17/2/18.
//  Copyright © 2017年 吴道睿. All rights reserved.
//

package rlog

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"time"
)

type severity int32

const (
	infoLog severity = iota
	monitor
	numSeverity = 2
)

//日志标志串  分别对应于下面数组内容的首字母
const severityChar = "IM"

var severityName = []string{
	infoLog:    "INFO",
	monitor: "Monitor",
}

var b_severityName = [][]byte{
	infoLog:  []byte("INFO ]"),//[]byte("I] "),
	monitor:  []byte("M] "),
}

var severityStats = [numSeverity]*OutputStats{
	infoLog:    &Stats.Info,
	monitor: &Stats.Monitor,
}

//
type flushSyncWriter interface {
	Flush() error
	Sync() error
	io.Writer
}

// 把所有内容写到磁盘
func Flush() {
	logging.lockAndFlushAll()
}


var timeNow = time.Now

/*
设置每行日志头.

每行日志的格式:
	Lmmdd hh:mm:ss.uuuuuu threadid file:line] msg...
格式说明:
	L                日志级别的首字母 (如 'I' 代表 INFO)
	mm               月份 (不足2位补0; 如 '05' 表示5月)
	dd               日期 (不足2位补0)
	hh:mm:ss.uuuuuu  时分秒 微秒
	threadid         线程id
	file             文件名
	line             行号
	msg              要输出的信息
*/
func (l *loggingT) header(s severity) *buffer {
	// Lmmdd hh:mm:ss.uuuuuu threadid file:line]
	now := timeNow()
	if s > monitor {
		s = monitor // 防止越界
	}
	buf := l.getBuffer()
	//为了不影响 fprintf的速度手动实现，速度是fprintf的30多倍
	year, month, day := now.Date()
	hour, minute, second := now.Clock()
	
	buf.tmp[0] = '['
	buf.tmp[1] = ' '
	buf.nDigits(4,2,int64(year))
	buf.twoDigits(6, int(month))
	buf.twoDigits(8, day)
	buf.tmp[10] = ' '
	buf.twoDigits(11, hour)
	buf.tmp[13] = ':'
	buf.twoDigits(14, minute)
	buf.tmp[16] = ':'
	buf.twoDigits(17, second)
	buf.tmp[19] = '.'
	buf.nDigits(6, 20, int64(now.Nanosecond()/1000))
	buf.tmp[26] = ' '
	buf.Write(buf.tmp[:27])
    buf.Write(b_severityName[s])
	return buf
}


//输出一行日志
func (l *loggingT) println(s severity, args ...interface{}) {
	buf := l.header(s)
	fmt.Fprintln(buf, args...)
	l.output(s, buf)
}

//输出日志
func (l *loggingT) print(s severity, args ...interface{}) {
	buf := l.header(s)
	fmt.Fprint(buf, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf)
}

//按格式输出日志
func (l *loggingT) printf(s severity, format string, args ...interface{}) {
	buf := l.header(s)
	fmt.Fprintf(buf, format, args...)
	if buf.Bytes()[buf.Len()-1] != '\n' {
		buf.WriteByte('\n')
	}
	l.output(s, buf)
}



type syncBuffer struct {
	logger *loggingT
	*bufio.Writer
	file   *os.File
	sev    severity
	nbytes uint64 // 写入文件的字节数
    day    int
}

func (sb *syncBuffer) Sync() error {
	return sb.file.Sync()
}

//写数据   当数据文件的长度大于最大的文件的长度 会自动生成新的文件
func (sb *syncBuffer) Write(p []byte) (n int, err error) {
    //换天时新建文件
	day := time.Now().Day()
	if sb.nbytes+uint64(len(p)) >= MaxSize  || day != sb.day{
	    sb.day = day
		if err := sb.rotateFile(time.Now()); err != nil {
			sb.logger.exit(err)
		}
	}
	n, err = sb.Writer.Write(p)
	sb.nbytes += uint64(n)
	if err != nil {
		sb.logger.exit(err)
	}
	return
}

//写缓存的大小，满了后才会写入到文件中除非主动调用flash 256k
var bufferSize int= 256 * 1024

func (sb *syncBuffer) rotateFile(now time.Time) error {
	if sb.file != nil {
		sb.Flush()
		sb.file.Close()
	}
	var err error
	sb.file, _, err = create(severityName[sb.sev], now)
	sb.nbytes = 0
	if err != nil {
		return err
	}
	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)
	// 写文件头.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, "Log file created at: %s\n", now.Format("2006/01/02 15:04:05"))
	fmt.Fprintf(&buf, "Binary: Built with %s %s for %s/%s\n", runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintf(&buf, "Log line format: [ISW]mmdd hh:mm:ss.uuuuuu threadid ] msg\n")
	n, err := sb.file.Write(buf.Bytes())
	sb.nbytes += uint64(n)
	return err
}

//创建文件
func (l *loggingT) createFiles(sev severity) error {
	now := time.Now()
	// 根据 severity递减创建文件 如果存在则不建
	sb := &syncBuffer{
		logger: l,
		sev:    sev,
		day:    time.Now().Day(),
	}
	if err := sb.rotateFile(now); err != nil {
		return err
	}
	l.file[sev] = sb
	return nil
}

//////////////////////////////////////////////////////////////////////
type Verbose bool

//外部调用函数
func V(level Level) Verbose {
	if logging.verbosity.get() >= level {
		return Verbose(true)
	}
	return Verbose(false)
}

//－－－－－－－－－－－－－－－
func (v Verbose) Info(args ...interface{}) {
	if v {
		logging.print(infoLog, args...)
	}
}

/////////////////////////////////////////////////////////////////////////
func Info(args ...interface{}) {
	logging.print(infoLog, args...)
}
