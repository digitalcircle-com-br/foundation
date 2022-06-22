package core

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"syscall"
	"time"

	_ "github.com/breml/rootcerts"
	"github.com/digitalcircle-com-br/buildinfo"
	"github.com/google/uuid"
)

func init() {
	Log("Initiating foundation v0.0.11")
}

func Ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second*120)
}

func IsDocker() bool {
	_, err := os.Stat("/.dockerenv")
	return err == nil
}

var svcName = "foundation"
var svcId = NewUUID()

var sigCh = make(chan os.Signal)

var onStop = func() {

}

func SvcName() string {
	return svcName
}

func Init(s string) {
	svcName = s
	signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		<-sigCh

		os.Exit(0)
	}()
	if IsDocker() {
		Log("Initiating Container for: %s", svcName)
	} else {
		Log("Initiating Service: %s", svcName)
	}
	wd, _ := os.Getwd()
	abswd, _ := filepath.Abs(wd)
	Log("Running from %s", abswd)
	Log(buildinfo.String())
}

func Log(s string, p ...interface{}) {

	log.Printf(fmt.Sprintf("LOG [%s] - %s", svcName, s), p...)
}

func Warn(s string, p ...interface{}) {
	bs := debug.Stack()
	log.Printf(fmt.Sprintf("WARN [%s] - %s\n\t%s", svcName, s, string(bs)), p...)
	//log.Printf(fmt.Sprintf("WARN [%s] - %s", svcName, s), p...)

}

func Debug(s string, p ...interface{}) {
	log.Printf(fmt.Sprintf("DBG [%s] - %s", svcName, s), p...)
}

func Fatal(s ...interface{}) {
	log.Fatal(s...)
}

func Err(err error) {
	if err != nil {
		Log("Error: %s", err.Error())
	}
}

func NewUUID() string {
	return uuid.NewString()
}

func LogWriter() io.Writer {
	return log.Default().Writer()
}
