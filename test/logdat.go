//
//  logdat.go
//  rlog
//
//  Created by 吴道睿 on 17/2/16.
//  Copyright © 2017年 吴道睿. All rights reserved.
//
package main

import(
    "../../rlog"
	"time"
)

func main(){
	rlog.LogInit(3,"./log",1800000,1)
	rlog.StartMonitorLog("./log")
	
	logstr := `,"tid":"b8F2H2GEALtJP0YT0","org_type":"Mobile","org_name":"HE_10086","method":"rush","error_message":"","error_type":"","status":"success","rt":"1234"`
	
	dat := map[string]string{"tid":"b8F2H2GEALtJP0YT0","org_type":"Mobile","org_name":"HE_10086","method":"rush","error_message":"","error_type":"","status":"success","rt":"1234"}
	
	for  i:= 0;i<10;i++{
	   time.Sleep(1000)
       rlog.M(1).Log(logstr)
	   rlog.M(1).Info(dat)
	}
	
	for  i:= 0;i<10;i++{
		rlog.V(1).Info(`testlogbcisdhcvd,ssdhbchd`,10,dat)
	}
}
