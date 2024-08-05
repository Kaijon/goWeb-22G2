package main

import (
	"getac/goWeb/utils"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
)

func main() {
	numProcs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(numProcs)
	Log.Infoln("GOMAXPROCS set to:", numProcs)
	workDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		Log.Fatal(err)
	}
	Log.Infoln("WORKDIR:", workDir)
	out, err := exec.Command("/usr/sbin/iptables -F").Output()
	if err != nil {
		log.Println(err)
	}
	Log.Infoln("Run iptables -F", out)
	//log.SetFlags(0)
	//log.SetOutput(new(logWriter))
	LogInit()
	//ListSsdp.Start()
	//go StartSsdpLoop()
	go utils.StartMqttInLoop()
	go utils.StartMqttExLoop()
	go serveHTTP()
	go serveStreams()
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		Log.Println(sig)
		done <- true
	}()
	Log.Println("Server Start Awaiting Signal")
	<-done
	Log.Println("Exiting")
}
