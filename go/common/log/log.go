package log

// package log
// use log must invoke Init() first,or it will panic
//usage:
// log.Init(nil)
// defer log.Flush()

import (
	"fmt"
	sl "github.com/cihub/seelog"
)

func init(){
	Init(nil)
}
var logWorker sl.LoggerInterface

type Level int

const (
	debug Level = 1
	warn        = 2
	error       = 3
)

type Config struct {
	FileName string
	MaxRoll  string
	MinLevel Level
	CanStd   bool
}

func Init(cfg *Config) {
	var f, m, l, s string
	if cfg == nil {
		f = "default.log"
		m = "1"
		l = getLevels(debug)
	} else {
		f = cfg.FileName
		m = cfg.MaxRoll
		l = getLevels(cfg.MinLevel)
		if cfg.CanStd {
			s = "<console/>"
		}
	}
	config := `
<seelog type="sync" levels="%s" >
    <outputs formatid="main">
        <rollingfile type="date" filename="%s" datepattern="02.01.2006" maxrolls="%s"/>
        %s
    </outputs>
    <formats>
        <format id="main" format="[%%Date %%Time] [%%LEVEL] %%Func(): %%Msg%%n"/>
    </formats>
</seelog>`
	config = fmt.Sprintf(config, l, f, m, s)
	logger, err := sl.LoggerFromConfigAsBytes([]byte( config ))
	if err != nil {
		panic("log init failed")
	}
	loggerErr := sl.ReplaceLogger(logger)
	if loggerErr != nil {
		panic("log init failed")
	}
	logWorker = logger
}

func Debugf(format string, params ...interface{}) {
	logWorker.Debugf(format, params)
}
func Warnf(format string, params ...interface{}) {
	logWorker.Warnf(format, params)
}
func Errorf(format string, params ...interface{}) {
	logWorker.Errorf(format, params)
}
func Debug(param interface{}) {
	logWorker.Debug(param)
}
func Warn(param interface{}) {
	logWorker.Warn( param)
}
func Error( param interface{}) {
	logWorker.Error(param)
}
func Flush() {
	logWorker.Flush()
}

func getLevels(l Level) string {
	switch l {
	case debug:
		return "debug,warn,error"
	case warn:
		return "warn,error"
	case error:
		return "error"
	default:
		return "error"
	}
}
