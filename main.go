package main

import (
	"github.com/codegangsta/cli"
	"github.com/project-5/rabbitmq-cli-consumer/command"
	"github.com/project-5/rabbitmq-cli-consumer/config"
	"github.com/project-5/rabbitmq-cli-consumer/consumer"
	"io"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "rabbitmq-cli-consumer"
	app.Usage = "Consume RabbitMQ easily to any cli program"
	app.Author = "Richard van den Brand"
	app.Email = "richard@vandenbrand.org"
	app.Version = "1.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "executable, e",
			Usage: "Location of executable",
		},
		cli.StringFlag{
			Name:  "configuration, c",
			Usage: "Location of configuration file",
		},
		cli.StringFlag{
			Name:  "host",
			Usage: "Rabbitmq host. This flug overrdide configuration option.",
		},
		cli.StringFlag{
			Name:  "port",
			Usage: "Rabbitmq port. This flug overrdide configuration option.",
		},
		cli.StringFlag{
			Name:  "user",
			Usage: "Rabbitmq user. This flug overrdide configuration option.",
		},
		cli.StringFlag{
			Name:  "password",
			Usage: "Rabbitmq password. This flug overrdide configuration option.",
		},
	}
	app.Action = func(c *cli.Context) {
		if c.String("configuration") == "" && c.String("executable") == "" {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		logger := log.New(os.Stderr, "", log.Ldate|log.Ltime)
		cfg, err := config.LoadAndParse(c.String("configuration"))

		if err != nil {
			logger.Fatalf("Failed parsing configuration: %s\n", err)
		}
		
		// override config parameters from flugs
		host := c.String("host");
		port := c.String("port");
		user := c.String("user");
		password := c.String("password");
		
		if "" != host {
			cfg.RabbitMq.Host = host;
		}
		if "" != port {
			cfg.RabbitMq.Port = port;
		}
		if "" != user {
			cfg.RabbitMq.Username = user;
		}
		if "" != password {
			cfg.RabbitMq.Password = password;
		}

		errLogger := log.New(io.MultiWriter(os.Stdout, os.Stderr), "", log.Ldate|log.Ltime)
		infLogger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

		factory := command.Factory(c.String("executable"))

		client, err := consumer.New(cfg, factory, errLogger, infLogger)
		if err != nil {
			errLogger.Fatalf("Failed creating consumer: %s", err)
		}

		client.Consume()
	}

	app.Run(os.Args)
}
