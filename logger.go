package rlog

type Logger interface{
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
	Info(v ...interface{})
	InfoDepth(depth int, v ...interface{}) 
	Infoln(v ...interface{})
	Infof(format string, v ...interface{})
	Warning(v ...interface{})
	WarningDepth(depth int, v ...interface{})
	Warningln(v ...interface{})
	Warningf(format string, v ...interface{})
	Error(v ...interface{}) 
	ErrorDepth(depth int, v ...interface{}) 
	Errorln(v ...interface{})
	Errorf(format string, v ...interface{}) 
	Fatal(v ...interface{})
	FatalDepth(depth int, v ...interface{}) 
	Fatalln(v ...interface{})
	Fatalf(format string, v ...interface{}) 
}

type Rlogger struct{}

func NewLogger()*Rlogger{
	return new(Rlogger)
}

func (r *Rlogger)Print(v ...interface{}){
    logging.print(infoLog, v...)
}

func (r *Rlogger)Printf(format string, v ...interface{}){
    logging.print(infoLog, v...)
}

func (r *Rlogger)Println(v ...interface{}){
    logging.print(infoLog, v...)
}

func (r *Rlogger)Info(v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)InfoDepth(depth int, v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)Infoln(v ...interface{}) {
	logging.println(infoLog, v...)
}

func (r *Rlogger)Infof(format string, v ...interface{}) {
	logging.printf(infoLog, format, v...)
}

func (r *Rlogger)Warning(v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)WarningDepth(depth int, v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)Warningln(v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)Warningf(format string, v ...interface{}) {
	logging.print(infoLog, v...)
}


func (r *Rlogger)Error(v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)ErrorDepth(depth int, v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)Errorln(v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)Errorf(format string, v ...interface{}) {
	logging.print(infoLog, v...)
}


func (r *Rlogger)Fatal(v ...interface{}) {
	logging.print(infoLog, v...)
}


func (r *Rlogger)FatalDepth(depth int, v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)Fatalln(v ...interface{}) {
	logging.print(infoLog, v...)
}

func (r *Rlogger)Fatalf(format string, v ...interface{}) {
	logging.print(infoLog, v...)
}
