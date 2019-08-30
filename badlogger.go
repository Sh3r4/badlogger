package badlogger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/fatih/color"
)

var (
	// singleton
	one sync.Once

	// logger handles
	h Handles = crappyInit()

	// colours
	Green   = color.New(color.FgGreen, color.Bold)
	Yellow  = color.New(color.FgYellow, color.Bold)
	Red     = color.New(color.FgRed, color.Bold)
	Magenta = color.New(color.FgMagenta, color.Bold)
	Blue    = color.New(color.FgBlue, color.Bold)

	// prefixes
	pDebug string = "[-] "
	pLog   string = "[*] "
	pWarn  string = "[!] "
	pErr   string = "[ERROR] "
)

type (
	// Handles wraps a number of log handles
	Handles struct {
		// logHandles
		debug    *log.Logger
		debugNP  *log.Logger
		log      *log.Logger
		logNP    *log.Logger
		warn     *log.Logger
		error    *log.Logger
		fatality *log.Logger
	}

	// CustomConfig represents a structure for setting the logging up using pre-built log handles
	CustomConfig struct {
		LogHandles struct {
			Debug         *log.Logger
			DebugNoPrefix *log.Logger
			Log           *log.Logger
			LogNoPrefix   *log.Logger
			Warn          *log.Logger
			Error         *log.Logger
			Fatal         *log.Logger
		}
		Prefixes struct {
			Debug string
			Log   string
			Warn  string
			Err   string
		}
	}
)

func crappyInit() Handles {
	h.debug = setupLoggerWithDatesAndTimes(ioutil.Discard, Green, pDebug)
	h.log = setupLoggerWithDatesAndTimes(os.Stdout, Green, pLog)
	h.debugNP = setupLogger(ioutil.Discard, Green, "")
	h.logNP = setupLogger(os.Stdout, Green, "")
	h.warn = setupLoggerWithDatesAndTimes(os.Stdout, Magenta, pWarn)
	h.error = setupLoggerWithDatesAndTimes(os.Stderr, Red, pErr)
	return h
}

// NewCustomLogging applies the CustomConfig to the single logging instance
func NewCustomLogging(config CustomConfig) {
	one.Do(func() {
		h.debug = config.LogHandles.Debug
		h.debugNP = config.LogHandles.DebugNoPrefix
		h.log = config.LogHandles.Log
		h.logNP = config.LogHandles.LogNoPrefix
		h.warn = config.LogHandles.Warn
		h.error = config.LogHandles.Error
		h.fatality = config.LogHandles.Fatal
		pDebug = config.Prefixes.Debug
		pLog = config.Prefixes.Log
		pWarn = config.Prefixes.Warn
		pErr = config.Prefixes.Err
	})
}

// NewBasicLogger sets up logging based on the log leve
// fatality will ALWAYS output logs
// dateStamps = true will add dates to all log output
// timeStamps = true will add times to all log output
// logLevel >=4 outputs all logs; =3 outputs log,warn,error; =2 outputs warn,error; =1 outputs error, =0 suppresses all
func NewBasicLogger(logLevel int, dateStamps bool, timeStamps bool) {

	one.Do(func() {

		// highest logging level
		if logLevel >= 4 {
			setupLoggers(
				os.Stdout,
				os.Stdout,
				os.Stdout,
				os.Stderr,
				dateStamps,
				timeStamps,
			)
		} else if logLevel == 3 {
			setupLoggers(
				ioutil.Discard,
				os.Stdout,
				os.Stdout,
				os.Stderr,
				dateStamps,
				timeStamps,
			)
		} else if logLevel == 2 {
			setupLoggers(
				ioutil.Discard,
				ioutil.Discard,
				os.Stdout,
				os.Stderr,
				dateStamps,
				timeStamps,
			)
		} else if logLevel == 1 {
			setupLoggers(
				ioutil.Discard,
				ioutil.Discard,
				ioutil.Discard,
				os.Stderr,
				dateStamps,
				timeStamps,
			)
		} else {
			setupLoggers(
				ioutil.Discard,
				ioutil.Discard,
				ioutil.Discard,
				ioutil.Discard,
				dateStamps,
				timeStamps,
			)
		}
	})
}

func Debug(s string) {
	h.debug.Println(s)
}

func DebugNP(s string) {
	h.debugNP.Println(s)
}

func Log(s string) {
	h.log.Println(s)
}

func LogNP(s string) {
	h.logNP.Println(s)
}

func Warn(s string) {
	h.warn.Println(s)
}

func Error(s string, err error) {
	if len(s) == 0 {
		h.error.Println(err)
	} else {
		h.error.Println(fmt.Sprintf("%s | %s", s, err))
	}
}

func Fatal(s string, err error) {
	if len(s) == 0 {
		h.fatality.Fatalln(err)
	} else {
		h.fatality.Fatalln(fmt.Sprintf("%s | %s", s, err))
	}
}

// ChangePrefixes sets up the log prefixes
// setting params to "" or "default" will use default values
func ChangePrefixes(debug string, log string, warn string, err string) {
	if debug == "default" || len(debug) == 0 {
		pDebug = "[-] "
	} else {
		pDebug = debug
	}
	if log == "default" || len(log) == 0 {
		pLog = "[*] "
	} else {
		pLog = log
	}
	if warn == "default" || len(warn) == 0 {
		pWarn = "[!] "
	} else {
		pWarn = warn
	}
	if err == "default" || len(err) == 0 {
		pErr = "[ERROR] "
	} else {
		pErr = err
	}
}

// ErrCheck determines whether an error occurred
// returns true if an error, false otherwise
// user ErrCheckLog to log the error
// use ErrCheckFatal to log and terminate
func ErrCheck(err error) bool {
	if err != nil {
		return true
	}
	return false
}

// ErrCheckLog will log an error if it exists and allow the program to continue
// returns true if an error was encountered
// Use ErrCheckFatal to log a fatal error
func ErrCheckLog(custom string, err error) bool {
	if err != nil {
		Error(custom, err)
		return true
	}
	return false
}

// ErrCheckFatal will log an error if it exists exit
func ErrCheckFatal(custom string, err error) {
	if err != nil {
		Fatal(custom, err)
	}
}

func setupLoggers(DebugWriter io.Writer, LogWriter io.Writer, warnWriter io.Writer, errorWriter io.Writer, dateStamps bool, timeStamps bool) {

	if dateStamps && timeStamps {
		h.debug = setupLoggerWithDatesAndTimes(DebugWriter, Green, pDebug)
		h.log = setupLoggerWithDatesAndTimes(LogWriter, Green, pLog)
		h.debugNP = setupLogger(DebugWriter, Green, "")
		h.logNP = setupLogger(LogWriter, Green, "")
		h.warn = setupLoggerWithDatesAndTimes(warnWriter, Magenta, pWarn)
		h.error = setupLoggerWithDatesAndTimes(errorWriter, Red, pErr)
	} else if dateStamps && !timeStamps {
		h.debug = setupLoggerWithDates(DebugWriter, Green, pDebug)
		h.log = setupLoggerWithDates(LogWriter, Green, pLog)
		h.debugNP = setupLogger(DebugWriter, Green, "")
		h.logNP = setupLogger(LogWriter, Green, "")
		h.warn = setupLoggerWithDates(warnWriter, Magenta, pWarn)
		h.error = setupLoggerWithDates(errorWriter, Red, pErr)
	} else if timeStamps && !dateStamps {
		h.debug = setupLoggerWithTimes(DebugWriter, Green, pDebug)
		h.log = setupLoggerWithTimes(LogWriter, Green, pLog)
		h.debugNP = setupLogger(DebugWriter, Green, "")
		h.logNP = setupLogger(LogWriter, Green, "")
		h.warn = setupLoggerWithTimes(warnWriter, Magenta, pWarn)
		h.error = setupLoggerWithTimes(errorWriter, Red, pErr)
	} else {
		h.debug = setupLogger(DebugWriter, Green, pDebug)
		h.log = setupLogger(LogWriter, Green, pLog)
		h.debugNP = setupLogger(DebugWriter, Green, "")
		h.logNP = setupLogger(LogWriter, Green, "")
		h.warn = setupLogger(warnWriter, Magenta, pWarn)
		h.error = setupLogger(errorWriter, Red, pErr)
	}
}

func setupLoggerWithDatesAndTimes(outputWriter io.Writer, colouriser *color.Color, prefixText string) (logger *log.Logger) {
	return log.New(outputWriter, colouriser.Sprintf(prefixText), log.Ldate|log.Ltime)
}

func setupLoggerWithDates(outputWriter io.Writer, colouriser *color.Color, prefixText string) (logger *log.Logger) {
	return log.New(outputWriter, colouriser.Sprintf(prefixText), log.Ldate)
}

func setupLoggerWithTimes(outputWriter io.Writer, colouriser *color.Color, prefixText string) (logger *log.Logger) {
	return log.New(outputWriter, colouriser.Sprintf(prefixText), log.Ltime)
}

func setupLogger(outputWriter io.Writer, colouriser *color.Color, prefixText string) (logger *log.Logger) {
	return log.New(outputWriter, colouriser.Sprintf(prefixText), 0)
}
