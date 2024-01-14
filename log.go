package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

const reset = "\033[0m"
const red = "\033[31m"
const green = "\033[32m"
const yellow = "\033[33m"
const blue = "\033[34m"
const gray = "\033[37m"

const grayBg = "\033[47m"
const blackBold = "\033[1;30m"

func writeStdErr(msg string) {
	lg := log.New(os.Stderr, "", 0)
	lg.Println(msg)
}

func getCaller() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	return runtime.FuncForPC(pc).Name()
}

func renderCaller(caller string) string {
	return gray + "[" + caller + "] " + reset
}

func renderTrace() string {
	var trace string
	for i := 2; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		if i == 2 {
			trace += gray + "\n  ↳ ⚙ TRACE: " + reset
			trace += fmt.Sprintf("%s:%d %s\n", file, line, runtime.FuncForPC(pc).Name())
			continue
		}

		trace += fmt.Sprintf("             %s:%d %s\n", file, line, runtime.FuncForPC(pc).Name())
	}

	return trace
}

type Logger struct {
	ModuleName string
	ShowDebug  bool
	ShowCaller bool
	ShowTrace  bool
	// Info, Warn, Error, Fatal, Debugのみ
	ShowTime bool
}

func (l *Logger) renderModuleName() string {
	if l.ModuleName != "" {
		return grayBg + blackBold + " " + l.ModuleName + " " + reset + " "
	}

	return ""
}

func (l *Logger) Info(text ...string) {
	out := l.renderModuleName()
	out += green + "ℹ INFO: " + reset
	if l.ShowTime {
		out += time.Now().Format("2006-01-02 15:04:05") + " "
	}

	out += strings.Join(text, " ")

	fmt.Println(out)
}

func (l *Logger) ProgressInfo(text ...string) {
	out := gray + "∴ INIT: " + reset
	out += strings.Join(text, " ")

	fmt.Println(out)
}

func (l *Logger) ProgressOk() {
	fmt.Println(green + "  ↳ ✔ OK!" + reset)
}

func (l *Logger) Error(text ...string) {
	out := l.renderModuleName()
	out += red + "✘ ERROR: " + reset
	if l.ShowTime {
		out += time.Now().Format("2006-01-02 15:04:05") + " "
	}

	if l.ShowCaller {
		out += renderCaller(getCaller())
	}

	out += strings.Join(text, " ")

	writeStdErr(out)
}

func (l *Logger) Fatal(text ...string) {
	out := l.renderModuleName()
	out += red + "💥 FATAL: " + reset
	if l.ShowTime {
		out += time.Now().Format("2006-01-02 15:04:05") + " "
	}

	if l.ShowCaller {
		out += renderCaller(getCaller())
	}

	out += strings.Join(text, " ")
	if l.ShowTrace {
		out += renderTrace()
	}

	writeStdErr(out)
}

func (l *Logger) ErrorWithDetail(text string, err error) {
	fmt.Fprintln(os.Stderr, red+"✘ ERROR: "+reset+text)
	fmt.Fprintln(os.Stderr, gray+"  ↳ ⚙ DETAIL: "+err.Error())
}

func (l *Logger) FatalWithDetail(text string, err error) {
	fmt.Fprintln(os.Stderr, red+"💥 FATAL: "+text+reset)
	fmt.Fprintln(os.Stderr, gray+"   ↳ ⚙ DETAIL: "+err.Error())
}

func (l *Logger) Warn(text ...string) {
	out := l.renderModuleName()
	out += yellow + "⚠ WARNING: " + reset
	if l.ShowTime {
		out += time.Now().Format("2006-01-02 15:04:05") + " "
	}

	out += strings.Join(text, " ")

	fmt.Println(out)
}

func (l *Logger) Debug(text ...string) {
	if l.ShowDebug {
		out := l.renderModuleName()
		out += blue + "⚙ DEBUG: " + reset
		if l.ShowTime {
			out += time.Now().Format("2006-01-02 15:04:05") + " "
		}

		out += strings.Join(text, " ")
		fmt.Println(out)
	}
}
