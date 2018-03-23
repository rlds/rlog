//
//  logbuf.go
//  rlog
//
//  Created by 吴道睿 on 17/2/18.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package rlog

import (
		"bytes"
)

const digits = "0123456789"

type buffer struct {
	bytes.Buffer
	tmp  [64]byte // 临时存储
	next *buffer
}


//数字转为2位字符 不足则在前面补0
func (buf *buffer) twoDigits(i, d int) {
	buf.tmp[i+1] = digits[d%10]
	d /= 10
	buf.tmp[i] = digits[d%10]
}

//数字转换为指定长度的字符串不足位数前补0
//需要考虑32位环境下的情况
func (buf *buffer) nDigits(n, i int, d int64) {
	for j := n - 1; j >= 0; j-- {
		buf.tmp[i+j] = digits[d%10]
		d /= 10
	}
}

//数字转换字符并复制至指定的i位后面
func (buf *buffer) someDigits(i, d int) int {
	j := len(buf.tmp)
	for {
		j--
		buf.tmp[j] = digits[d%10]
		d /= 10
		if d == 0 {
			break
		}
	}
	return copy(buf.tmp[i:], buf.tmp[j:])
}


