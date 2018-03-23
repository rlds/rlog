### 使用说明：

    首先需要下载此包至 GOPATH 路径下 
    
	在需要使用的main包中
    
	import(
	   "rlog"
	)
	
	在使用之前需要调用


	rlog.LogInit(conf.LogLevel,conf.LogDir,conf.MaxLogLen_m,conf.CacheLen)
	参数说明：
	conf.LogLevel    （数值 日志输出级别 大于0的数值）
	conf.LogDir      （日志最终要存放在那个文件夹下可以是绝对路径或相对路径）
	conf.MaxLogLen_m （日志文件的最大尺寸，当日志大于这一数值时兴建新文件）
	conf.CacheLen    （日志在内存中的缓存大小，即日志先存在内存中超过此值后写入文件，不需要缓存此值可以设为0）
	
	N 为数值 指定的日志级别 小于 conf.LogLevel的日志才会输出

### //文本格式日志输出：
	rlog.V(N).Info(...)

### //json格式日志输出
    //在启动程序开始时加入下面的启动代码
    rlog.StartMonitorLog("./") 
    
    //注意Info Log 会附加timestep字段 时间精确到微秒值 位于json串第一个字段
	rlog.M(N).Info(map[string]string)
    rlog.M(N).Log(dat string)  dat 格式:  ,"key1":"val1","key2":"val2"
    
    //任意结构体类型输出至日志 不会附加字符，每个结构体内容转为一行
    rlog.M(N).ILog(obj interface{})
   
### //日志自动清理
	rlog.SetRemoveLogFileDayNum( num )
    num > 0 有效 num 天数的数值
    num <=0 则表示不清理日志
    清理规则：删除最近日志日期 前推n天之前的日志，日志文件需包含yyyyMMdd格式日期字符串
    ⚠️⚠️ 可以不主动设置，默认清理最后一个日志文件9天以前的日志
    ⚠️⚠️ 有后期手动更新使得日志文件时间与文件名时间不符的日志文件不删除

