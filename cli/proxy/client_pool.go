package proxy

import (
	"errors"
	"log"
	"net/http"
	"sync"
	"time"
)

type httpResult struct {
	Res *http.Response
	Err error
}

type message struct {
	req *http.Request
	ch  chan *httpResult
}

type ClientPool struct {
	clients []*http.Client
	queue   chan message

	wg         sync.WaitGroup
	currClient int

	ticker     chan struct{}
	tickerStop chan struct{}
}

func (this *ClientPool) tick(rateLimit time.Duration) {
	defer close(this.ticker)
	var wg sync.WaitGroup

	for i := 0; i < len(this.clients); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				this.ticker <- struct{}{}

				select {
				case <-this.tickerStop:
					return
				case <-time.After(rateLimit):
					break
				}
			}
		}()
	}

	wg.Wait()
}

func (this *ClientPool) closeTicker() {
	close(this.tickerStop)
	for range this.ticker {
	}
}

func (this *ClientPool) start() {
	log.Println("client pool: started")
	defer this.wg.Done()

	for m := range this.queue {
		<-this.ticker
		c := this.client()
		res, err := c.Do(m.req)
		result := &httpResult{Res: res, Err: err}
		m.ch <- result
		close(m.ch)
	}
}

func (this *ClientPool) client() *http.Client {
	ret := this.clients[this.currClient]
	this.currClient++
	this.currClient %= len(this.clients)
	return ret
}

func (this *ClientPool) Shutdown() <-chan struct{} {
	log.Println("client pool: shutting down...")
	close(this.queue)
	stopped := make(chan struct{})

	go func() {
		this.wg.Wait()
		this.closeTicker()
		close(stopped)
		log.Println("client pool: shutted down")
	}()

	return stopped
}

func (this *ClientPool) Do(req *http.Request) (*http.Response, error) {
	m := message{
		req: req,
		ch:  make(chan *httpResult),
	}

	this.queue <- m
	res := <-m.ch
	<-m.ch

	return res.Res, res.Err
}

// rateLimit represents minimum rest time needed between two requests
// on a single proxy client.
func NewClientPool(
	rateLimit time.Duration, clients ...*http.Client,
) (*ClientPool, error) {
	if len(clients) == 0 {
		return nil, errors.New("empty client list")
	}

	m := &ClientPool{
		// concurrency is ok between proxies, just rate limit
		clients: clients,
		queue:   make(chan message, len(clients)),

		wg:         sync.WaitGroup{},
		currClient: 0,

		ticker:     make(chan struct{}, len(clients)),
		tickerStop: make(chan struct{}),
	}

	m.wg.Add(1)
	go m.tick(rateLimit)
	go m.start()

	return m, nil
}
