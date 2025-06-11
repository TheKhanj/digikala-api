package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/thekhanj/digikala-api/cli/config"
	"github.com/thekhanj/digikala-api/cli/fetch"
	"github.com/thekhanj/digikala-api/cli/proxy"
)

type Cli struct {
	config *config.Config
}

func (this *Cli) Run() Exit {
	var flags *flag.FlagSet
	if len(os.Args) == 0 {
		flags = flag.NewFlagSet("", flag.ExitOnError)
	} else {
		flags = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}

	configPath := flags.String("c", "config.json", "path to config file")
	flags.Parse(os.Args[1:])

	c, err := config.ReadConfig(*configPath)
	if err != nil {
		return NewExit(err, CODE_READING_CONFIG_ERROR)
	}

	this.config = c
	return this.handleCommand(flags.Args())
}

func (this *Cli) handleCommand(args []string) Exit {
	if len(args) == 0 {
		return NewHelpfulExit(
			errors.New("missing command"), CODE_BAD_INVOCATION_ERROR,
		)
	}

	cmd := args[0]

	switch cmd {
	case "fetch":
		return this.handleFetch(args[1:])
	case "gen-schema":
		return this.handleGenSchema(args[1:])
	default:
		return NewHelpfulExit(
			fmt.Errorf("invalid command: %s", cmd),
			CODE_BAD_INVOCATION_ERROR,
		)
	}
}

func (this *Cli) handleFetch(args []string) Exit {
	msg := `Try:
  digi fetch products`

	if len(args) == 0 {
		return NewMsgExit(
			errors.New("missing subcommand for fetch"),
			CODE_BAD_INVOCATION_ERROR, msg,
		)
	}

	subcmd := args[0]

	switch subcmd {
	case "products":
		return this.handleFetchProducts()
	default:
		return NewMsgExit(
			fmt.Errorf("invalid subcommand for fetch: %s", subcmd),
			CODE_BAD_INVOCATION_ERROR, msg,
		)
	}
}

func (this *Cli) handleGenSchema(args []string) Exit {
	return this.notImplemented()
}

func (this *Cli) notImplemented() Exit {
	return NewExit(
		errors.New("not implemented"), CODE_NOT_IMPLEMENTED_ERROR,
	)
}

func (this *Cli) handleFetchProducts() Exit {
	proxies, err := this.config.Api.Client.GetProxies()
	if err != nil {
		return NewExit(err, CODE_INVALID_CONFIG_ERROR)
	}

	clients, err := proxy.NewProxyClientList(proxies)
	if err != nil {
		return NewExit(err, CODE_INVALID_CONFIG_ERROR)
	}

	pool, err := proxy.NewClientPool(
		time.Duration(this.config.Api.Client.RateLimit), clients...,
	)
	if err != nil {
		return NewExit(err, CODE_INVALID_CONFIG_ERROR)
	}
	defer func() {
		stopped := pool.Shutdown()
		<-stopped
	}()

	urls, err := this.config.Api.Fetch.GetProductsApiUrls()
	if err != nil {
		return NewExit(err, CODE_INVALID_CONFIG_ERROR)
	}

	p, err := fetch.NewProducts(pool, urls, ".cache/fetch/products")
	if err != nil {
		return NewExit(err, CODE_GENERAL_ERROR)
	}

	err = p.Fetch()
	if err != nil {
		return NewExit(err, CODE_GENERAL_ERROR)
	}

	return NewExit(nil, CODE_SUCCESS)
}
