/*
© Copyright IBM Corporation 2017, 2022

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// chkmqready checks that MQ is ready for work, by checking if the MQ listener port is available
package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/ibm-messaging/mq-container/internal/ready"
	"github.com/ibm-messaging/mq-container/pkg/name"
)

func doMain() int {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()

	// Check if runmqserver has indicated that it's finished configuration
	r, err := ready.Check()
	if !r || err != nil {
		return 1
	}
	name, err := name.GetQueueManagerName()
	if err != nil {
		fmt.Println(err)
		return 1
	}

	// Check if the queue manager has a running listener
	if active, _ := ready.IsRunningAsActiveQM(ctx, name); active {
		conn, err := net.Dial("tcp", "127.0.0.1:1414")
		if err != nil {
			fmt.Println(err)
			return 1
		}
		err = conn.Close()
		if err != nil {
			fmt.Println(err)
		}
	} else if standby, _ := ready.IsRunningAsStandbyQM(ctx, name); standby {
		fmt.Printf("Detected queue manager running in standby mode")
		return 10
	} else if replica, _ := ready.IsRunningAsReplicaQM(ctx, name); replica {
		fmt.Printf("Detected queue manager running in replica mode")
		return 20
	} else {
		return 1
	}
	return 0
}

func main() {
	os.Exit(doMain())
}
