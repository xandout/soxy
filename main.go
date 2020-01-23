/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"github.com/xandout/soxy/client"
	"github.com/xandout/soxy/server"
	"os"
)

func main() {

	app := &cli.App{
		Name:  os.Args[0],
		Usage: "fight the loneliness!",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "soxy-url", Aliases: []string{"U"}, Usage: "https://soxy-daemon.com"},
			&cli.StringFlag{Name: "local", Aliases: []string{"L"}, Usage: "Which local port to listen on.\n\tExample: 3306 or 0.0.0.0:3306"},
			&cli.StringFlag{Name: "remote", Aliases: []string{"R"}, Usage: "Where should the daemon proxy traffic to?\n\tExample: mysql-service:3306"},
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:  "serve",
			Usage: "Start proxying traffic(server)",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "port", Aliases: []string{"p"}},
			},
			Action: server.Start,
		},
		{
			Name:  "proxy",
			Usage: "Start proxying traffic(client)",
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "soxy-url", Aliases: []string{"U"}, Usage: "ws://soxy-daemon.com:8080"},
				&cli.StringFlag{Name: "local", Aliases: []string{"L"}, Usage: "Which local port to listen on.\n\tExample: 3306 or 0.0.0.0:3306"},
				&cli.StringFlag{Name: "remote", Aliases: []string{"R"}, Usage: "Where should the daemon proxy traffic to?\n\tExample: mysql-service:3306"},
			},
			Action: client.Start,
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
