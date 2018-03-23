//
//  loglevel.go
//  rlog
//
//  Created by 吴道睿 on 17/2/18.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package rlog

import (
	"sync/atomic"
)

type Level int32

func (l *Level) get() Level {
	return Level(atomic.LoadInt32((*int32)(l)))
}

func (l *Level) set(val Level) {
	atomic.StoreInt32((*int32)(l), int32(val))
}


//用于外部动态更改日志输出级别
func SetLevel(v int32){
	logging.verbosity.set(Level(v))
	objLoging.verbosity.set(Level(v))
}

