package rlog

type Logger interface{
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
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