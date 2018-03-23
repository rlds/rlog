//
//  objlog.go
//  rlog
//
//  Created by 吴道睿 on 17/2/16.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package rlog

import(
	"bytes"
	"fmt"
	"os"
	"sync/atomic"
	"time"
	"encoding/json"
)

/*
   用于结构化输出日志
   输出格式为json格式 每个json一行
   单层json结构字符串类型
   此结构化日志为定制输出
*/

var objLoging loggingT

//
func (l *loggingT) output_obj(s severity, buf *buffer) {
	l.mu.Lock()
	data := buf.Bytes()
	if l.file[s] == nil {
		if err := l.createFiles_obj(s); err != nil {
			os.Stderr.Write(data) // 确保消息会输出.
			l.exit(err)
		}
	}
	l.file[monitor].Write(data)
	l.putBuffer(buf)
	l.file[monitor].Flush()
	l.mu.Unlock()
	if stats := severityStats[s]; stats != nil {
		atomic.AddInt64(&stats.lines, 1)
		atomic.AddInt64(&stats.bytes, int64(len(data)))
	}
}

func (l *loggingT) createFiles_obj(sev severity) error {
	now := time.Now()
	// 根据 severity递减创建文件 如果存在则不建
	// 只创建单独的文件
	if l.file[sev] == nil {
		sb := &objsyncBuffer{
			logger: l,
			sev:    sev,
			day:    time.Now().Day(),
		}
		if err := sb.rotateFile(now); err != nil {
			return err
		}
		l.file[sev] = sb
	}
	return nil
}

//时间输出毫秒值
func (l *loggingT) ts_header(s severity) *buffer {
	// Lmmdd hh:mm:ss.uuuuuu threadid file:line]
	now := timeNow()
	if s > monitor {
		s = monitor // 防止越界
	}
	buf := l.getBuffer()
	//为了不影响 fprintf的速度手动实现，速度是fprintf的30多倍
	copy(buf.tmp[0:13],[]byte(`{"timestamp":`))
	buf.nDigits(13,13,int64(now.UnixNano()/1000000))
	//{"timestamp":1487747214479
	buf.Write(buf.tmp[:26])
	return buf
}


func (l *loggingT) obj_println(s severity, args ...interface{}) {
	buf := l.ts_header(s) //l.getBuffer()
	fmt.Fprintln(buf, args...)
	l.output_obj(s, buf)
}


type ObjVerbose bool
type LogMap map[string]string

type ObjMsgVer struct{
	//lv Level
	l bool
	lm map[string]string
	msg string
	buf bytes.Buffer
}

func NewObjMsgVer(ok bool)(o *ObjMsgVer ){
	o = new (ObjMsgVer)
	if ok {
	    o.l  = true
	    o.lm = make(map[string]string)
		o.msg = ``
	}else{
		o.l  = false
	}
	return 
}

//外部调用函数
func M(level Level) * ObjMsgVer {
	if objLoging.verbosity.get() >= level {
		return NewObjMsgVer(true)
	}
	return NewObjMsgVer(false)
}

//
func (m *ObjMsgVer)Info(logmap map[string]string){
	if  m.l && logmap != nil {
		for k , v := range logmap {
		    m.msg += `,"`+ k +`":"` + v + `"`
		}
		m.msg += `}`
		objLoging.obj_println(monitor,m.msg)
	}
	return 
}

/*
   格式: ,"key1":"val1","key2":"val2"
*/
func (m *ObjMsgVer)Log(logstr string){
	if  m.l  {
		m.msg += logstr + `}`
		objLoging.obj_println(monitor,m.msg)
	}
	return 
}

//结构化日志 传入结构体，将被转成字符串 每行一条结构体的内容
func (m *ObjMsgVer)ILog(logobj interface{}){
	if  m.l  {
		dat,err := json.Marshal(logobj)
		if err == nil {
			dat = append(dat,'\n')
			buf := objLoging.getBuffer()
			buf.Write(dat)
			objLoging.output_obj(monitor,buf)
		}
	}
	return 
}

