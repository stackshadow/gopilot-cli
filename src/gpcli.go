/*
Copyright (C) 2019 by Martin Langlotz aka stackshadow

This file is part of gopilot, an rewrite of the copilot-project in go

gopilot is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 3 of this License

gopilot is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with gopilot.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"gopilot/clog"
	"gopilot/gbus"
	"gopilot/nodeName"
	"time"
)

var cmdGroup string
var cmdCommand string
var cmdPayload string
var listenToMessaged bool
var doPing bool

func main() {

	// ########################## Command line parse ##########################
	// core stuff
	clog.ParseCmdLine()
	mynodename.ParseCmdLine()
	flag.StringVar(&cmdGroup, "group", "", "The service")
	flag.StringVar(&cmdCommand, "cmd", "", "The command to the group")
	flag.StringVar(&cmdPayload, "payload", "", "The payload to the service")
	flag.BoolVar(&listenToMessaged, "listen", false, "Will listen to all messages on the socket")
	flag.BoolVar(&doPing, "ping", false, "Ping the server")
	flag.Parse()

	// ########################## Init ##########################
	clog.Init()
	mynodename.Init()

	if cmdCommand != "" {
		runCommand()
		return
	}

	if listenToMessaged == true {
		var listenbus gbus.Socketbus

		listenbus.Init()
		listenbus.Subscribe("","", "", func(message *gbus.Msg, group, command, payload string) {
			fmt.Printf("%+v\n", message)
		})
		listenbus.Connect(
			gbus.SocketFileName,
			gbus.Msg{NodeSource: "", GroupSource: ""},
		)
	}

	if doPing == true {
		runPing()
	}

}


func runPing() {
	var pingbus gbus.Socketbus
	pingbus.Init()
	remoteNodeName, _ := pingbus.Connect(
		gbus.SocketFileName,
		gbus.Msg{NodeSource: mynodename.NodeName, GroupSource: "cli"},
	)

	pingbus.Subscribe("",mynodename.NodeName, "cli", func(message *gbus.Msg, group, command, payload string) {
		fmt.Printf("%+v\n", message)
	})

	// create a new message
	pingbus.PublishPayload(
		mynodename.NodeName,
		remoteNodeName,
		"cli",
		"core",
		"ping",
		"",
	)

	time.Sleep(10 * time.Second)
}

func runCommand() {
	var cmdbus gbus.Socketbus
	cmdbus.Init()
	remoteNodeName, _ := cmdbus.Connect(
		gbus.SocketFileName,
		gbus.Msg{NodeSource: "", GroupSource: "cli"},
	)

	// create a new message
	cmdbus.PublishPayload(
		mynodename.NodeName,
		remoteNodeName,
		"cli",
		cmdGroup,
		cmdCommand,
		cmdPayload,
	)

	time.Sleep(10 * time.Second)
}
