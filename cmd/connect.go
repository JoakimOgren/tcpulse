/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
package cmd

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

// connectCmd represents the connect command
var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "connect connects to a port where 'serve' listens",
	RunE: func(cmd *cobra.Command, args []string) error {
		connections := cmd.Flags().Int32P("connections", "c", 10, "Number of connections to keep")
		rate := cmd.Flags().Int32P("rate", "r", 100, "connections throughput (/s)")
		duration := cmd.Flags().DurationP("duration", "d", 10*time.Second, "measurement period")

		addrport := cmd.Flags().Arg(0)

		return connect(addrport, &connectOptions{
			Connections: *connections,
			Rate:        *rate,
			Duration:    *duration,
		})
	},
}

func init() {
	rootCmd.AddCommand(connectCmd)
}

type connectOptions struct {
	Connections int32
	Rate        int32
	Duration    time.Duration
}

func connect(addrport string, opt *connectOptions) error {
	wg := &sync.WaitGroup{}
	var i int32
	for i = 0; i < opt.Connections; i++ {
		wg.Add(1)
		go func() {
			conn, err := net.Dial("tcp", addrport)
			if err != nil {
				log.Printf("could not dial %q: %s", addrport, err)
			}
			if _, err := conn.Write([]byte("Hello")); err != nil {
				log.Printf("could not write: %s\n", err)
			}

			timer := time.NewTimer(opt.Duration)
			<-timer.C

			if err := conn.Close(); err != nil {
				log.Printf("could not close: %s\n", err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}
