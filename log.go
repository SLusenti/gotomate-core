package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

// MyCaller returns the caller of the function that called it :)
func myCaller() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	return fmt.Sprintf("%s:%d", "github.com"+strings.Split(getFrame(2).File, "github.com")[1], getFrame(2).Line)
}

var isInit = false

func LoggerInit(loglevel string, logpath string) {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	file, err := os.Open(logpath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	log.SetOutput(file)

	// Only log the warning severity or above.
	switch strings.ToLower(loglevel) {
	case "panic":
		log.SetLevel(log.PanicLevel)
	case "Fatal":
		log.SetLevel(log.FatalLevel)
	case "Error":
		log.SetLevel(log.ErrorLevel)
	case "Warn":
		log.SetLevel(log.WarnLevel)
	case "Info":
		log.SetLevel(log.InfoLevel)
	case "Debug":
		log.SetLevel(log.DebugLevel)
	case "Trace":
		log.SetLevel(log.TraceLevel)
	default:
		fmt.Println("ERROR: unknown LogLevel " + loglevel)
		os.Exit(1)
	}

	isInit = true
}

//ONLY for broken code
func LogPanic(s ...interface{}) {
	fmt.Print("PANIC: " + myCaller() + " ")
	fmt.Println(s...)
	if isInit {
		log.Panic(s)
	}
	os.Exit(1)
}

//could not proceed
func LogFatal(s ...interface{}) {
	fmt.Print("FATAL: " + myCaller() + " ")
	fmt.Println(s...)
	if isInit {
		log.Fatal(s...)
	}
	os.Exit(1)
}

//some non-bloking errors
func LogError(s ...interface{}) {
	fmt.Print("ERROR: " + myCaller() + " ")
	fmt.Println(s...)
	if isInit {
		log.Error(s...)
	}
}

//highlighted info
func LogWarn(s ...interface{}) {
	fmt.Print(" WARN: " + myCaller() + " ")
	fmt.Println(s...)
	if isInit {
		log.Warn(s...)
	}
}

//normal info
func LogInfo(s ...interface{}) {
	fmt.Print(" INFO: " + myCaller() + " ")
	fmt.Println(s...)
	if isInit {
		log.Info(s...)
	}
}

//some code vars
func LogDebug(s ...interface{}) {
	fmt.Print("DEBUG: " + myCaller() + " ")
	fmt.Println(s...)
	if isInit {
		log.Debug(s...)
	}
}

//fmt.println("your code block")
func LogTrace(s ...interface{}) {
	fmt.Print("TRACE: " + myCaller() + " ")
	fmt.Println(s...)
	if isInit {
		log.Trace(s...)
	}
}
