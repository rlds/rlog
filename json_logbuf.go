//
//  objlogbuf.go
//  rlog
//
//  Created by 吴道睿 on 17/2/17.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package rlog

import(
    "os"
	"bufio"
	"time"
	"runtime"
	"fmt"
	"bytes"
)

type objsyncBuffer struct {
	logger *loggingT
	*bufio.Writer
	file   *os.File
	sev    severity
	nbytes uint64 // 写入文件的字节数
	day    int
}

func (sb *objsyncBuffer) Sync() error {
	return sb.file.Sync()
}

//写数据   当数据文件的长度大于最大的文件的长度 会自动生成新的文件
func (sb *objsyncBuffer) Write(p []byte) (n int, err error) {
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

func (sb *objsyncBuffer) rotateFile(now time.Time) error {
	if sb.file != nil {
		sb.Flush()
		sb.file.Close()
	}
	var err error
	sb.file, _, err = create_monitor(severityName[sb.sev], now)
	sb.nbytes = 0
	if err != nil {
		return err
	}
	
	sb.Writer = bufio.NewWriterSize(sb.file, bufferSize)
	
	// 写文件头.
	var buf bytes.Buffer
	fmt.Fprintf(&buf, `{"timestamp":%d,"sysinfo":"%s %s %s %s","timeformat":"ms","logformat":"json","logversion":"v2.0"}`+"\n", now.UnixNano()/1000000 ,runtime.Compiler, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	n, err := sb.file.Write(buf.Bytes())
	sb.nbytes += uint64(n)
	return err
}
