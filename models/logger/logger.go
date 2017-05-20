package logmodels

import (
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/ebitgo/ConsoleColor"
)

var Logger *logs.BeeLogger

type LogLevel struct {
	Debug    bool
	Info     bool
	Trace    bool
	Warning  bool
	Error    bool
	Critical bool
}

var LogLevelConfig = &LogLevel{
	Debug:    true,
	Info:     true,
	Trace:    true,
	Warning:  true,
	Error:    true,
	Critical: true,
}

func (this *LogLevel) Reset() {
	this.Debug = false
	this.Critical = false
	this.Error = false
	this.Info = false
	this.Trace = false
	this.Warning = false
}

func InitLogger() {
	logPath := beego.AppConfig.String("server_log_path")
	Logger = logs.NewLogger(10000)
	Logger.EnableFuncCallDepth(true) // 输出调用的文件名和文件行号
	//filename 保存的文件名
	//maxlines 每个文件保存的最大行数，默认值 1000000
	//maxsize 每个文件保存的最大尺寸，默认值是 1 << 28, //256 MB
	//daily 是否按照每天 logrotate，默认是 true
	//maxdays 文件最多保存多少天，默认保存 7 天
	//rotate 是否开启 logrotate，默认是 true
	//level 日志保存的时候的级别，默认是 Trace 级别
	Logger.SetLogger("file", `{"filename":"`+logPath+"server.log"+`"}`)
}

func SPrintInfo(title, format string, v ...interface{}) (ret string) {
	ret = fmt.Sprintf(" [ %s ]\r\n\t", title)
	if v == nil {
		ret += fmt.Sprintf(format)
	} else {
		ret += fmt.Sprintf(format, v)
	}
	if LogLevelConfig.Info {
		ConsoleColor.Println(ConsoleColor.C_WHITE, time.Now().Format("2006-01-02 15:04:05")+" [I] "+ret)
	}
	return
}

func SPrintError(title, format string, v ...interface{}) (ret string) {
	ret = fmt.Sprintf("[ %s ]\r\n\t", title)
	if v == nil {
		ret += fmt.Sprintf(format)
	} else {
		ret += fmt.Sprintf(format, v)
	}
	if LogLevelConfig.Error {
		ConsoleColor.Println(ConsoleColor.C_RED, time.Now().Format("2006-01-02 15:04:05")+" [E] "+ret)
	}
	return
}

func SPrintDebug(title, format string, v ...interface{}) (ret string) {
	ret = fmt.Sprintf("[ %s ]\r\n\t", title)
	if v == nil {
		ret += fmt.Sprintf(format)
	} else {
		ret += fmt.Sprintf(format, v)
	}
	if LogLevelConfig.Debug {
		ConsoleColor.Println(ConsoleColor.C_CYAN, time.Now().Format("2006-01-02 15:04:05")+" [D] "+ret)
	}
	return
}

func SPrintWarning(title, format string, v ...interface{}) (ret string) {
	ret = fmt.Sprintf("[ %s ]\r\n\t", title)
	if v == nil {
		ret += fmt.Sprintf(format)
	} else {
		ret += fmt.Sprintf(format, v)
	}
	if LogLevelConfig.Warning {
		ConsoleColor.Println(ConsoleColor.C_BLUE, time.Now().Format("2006-01-02 15:04:05")+" [W] "+ret)
	}
	return
}

func SPrintTrace(title, format string, v ...interface{}) (ret string) {
	ret = fmt.Sprintf("[ %s ]\r\n\t", title)
	if v == nil {
		ret += fmt.Sprintf(format)
	} else {
		ret += fmt.Sprintf(format, v)
	}
	if LogLevelConfig.Trace {
		ConsoleColor.Println(ConsoleColor.C_YELLOW, time.Now().Format("2006-01-02 15:04:05")+" [T] "+ret)
	}
	return
}

func SPrintCritical(title, format string, v ...interface{}) (ret string) {
	ret = fmt.Sprintf("[ %s ]\r\n\t", title)
	if v == nil {
		ret += fmt.Sprintf(format)
	} else {
		ret += fmt.Sprintf(format, v)
	}
	if LogLevelConfig.Critical {
		ConsoleColor.Println(ConsoleColor.C_WHITE, time.Now().Format("2006-01-02 15:04:05")+" [C] "+ret)
	}
	return
}
