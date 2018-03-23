//
//  logstat.go
//  rlog
//
//  Created by 吴道睿 on 17/2/18.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package rlog

import (
	"sync/atomic"
)

/*
    日志输出行数和字节数统计
*/

//输出的行数和字节数记录
type OutputStats struct {
	lines int64
	bytes int64
}

//返回已经输出的行数
func (s *OutputStats) Lines() int64 {
	return atomic.LoadInt64(&s.lines)
}

// 返回已经记录的字节数值
func (s *OutputStats) Bytes() int64 {
	return atomic.LoadInt64(&s.bytes)
}

// 
var Stats struct {
	Info, Monitor OutputStats
}
