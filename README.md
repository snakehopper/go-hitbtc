go-spiral
==========

Unofficial Spiral API implementation written on Golang.

Inspired by ```saniales/go-hitbtc```

This version implement V1 Spiral API.

## Import
	import "github.com/snakehopper/go-spiral"
	
## Usage

In order to use the client with go's default http client settings you can do:

~~~ go
package main

import (
	"fmt"
	"github.com/snakehopper/go-spiral"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func main() {
	// spiral client
	spiral := spiral.New(API_KEY, API_SECRET)

	// Get balances
	balances, err := spiral.GetBalances()
	fmt.Println(err, balances)
}
~~~

In order to use custom settings for the http client do:

~~~ go
package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/snakehopper/go-spiral"
)

const (
	API_KEY    = "YOUR_API_KEY"
	API_SECRET = "YOUR_API_SECRET"
)

func main() {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	// spiral client
	bc := spiral.NewWithCustomHttpClient(conf.spiral.ApiKey, conf.spiral.ApiSecret, httpClient)

	// Get balances
	balances, err := spiral.GetBalances()
	fmt.Println(err, balances)

	// Initialize websocket connection
	client, err := spiral.NewWSClient()
	if err != nil {
		handleError(err) // do something
	}
	defer client.Close()

	// Subscribe and handle
	tickerFeed, err := client.SubscribeTicker("ETHBTC")
	for {
		ticker := <-tickerFeed
		fmt.Println(ticker)
	}


}
~~~

See ["Examples" folder for more... examples](https://github.com/snakehopper/go-spiral/blob/master/examples/spiral.go)

# Projects using this library

- Golang Crypto Trading Bot: a framework to create trading bots easily and seamlessly (https://github.com/snakehopper/golang-crypto-trading-bot)
