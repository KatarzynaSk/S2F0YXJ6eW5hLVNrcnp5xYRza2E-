package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/jmoiron/sqlx"
)

type fetcher struct {
	client http.Client
	mux    sync.Mutex
	jobs   map[int]context.CancelFunc
	db     *sqlx.DB
}

func NewFetcher(client http.Client, db *sqlx.DB) *fetcher {
	return &fetcher{
		client: client,
		db:     db,
		jobs:   make(map[int]context.CancelFunc),
	}
}

func (f *fetcher) StartActive() error {
	active, err := selectAllActiveRequest(f.db)
	if err != nil {
		return err
	}

	for _, r := range active {
		f.StartJob(r)
	}
	
	return nil
}

func (f *fetcher) StopJob(id int) {
	f.mux.Lock()
	defer f.mux.Unlock()

	cancelJob, ok := f.jobs[id]
	if !ok {
		return
	}

	cancelJob()
	delete(f.jobs, id)
}

func (f *fetcher) StartJob(r request) {
	ctx, cancel := context.WithCancel(context.Background())

	f.mux.Lock()
	f.jobs[r.ID] = cancel
	f.mux.Unlock()

	go func() {
		for {
			ticker := time.NewTicker(time.Second * time.Duration(r.Interval))

			select {
			case <-ticker.C:
				fetchResult, err := f.makeRequest(r)
				if err != nil {
					log.Print(err.Error())
				}

				err = insertRequestResult(f.db, fetchResult)
				if err != nil {
					log.Print(err.Error())
				}

			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
}

func (f *fetcher) makeRequest(r request) (requestResult, error) {
	fetchResult := requestResult{RequestID: r.ID}

	start := time.Now()
	resp, err := f.client.Get(r.URL)
	fetchResult.Duration = time.Since(start).Seconds()

	if err != nil {
		return fetchResult, err
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fetchResult, err
	}

	contentString := string(content)
	fetchResult.Response = &contentString

	return fetchResult, nil
}
