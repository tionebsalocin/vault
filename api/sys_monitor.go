package api

import (
	"bufio"
	//"bytes"
	"context"
	//"encoding/json"
	//"fmt"
	//log "github.com/hashicorp/go-hclog"
	//"os"
	//"time"
)

//func (c *Sys) Monitor(loglevel string, logJSON bool, stopCh <-chan struct{}, q *QueryOptions) (chan string, error) {
// TODO: this function should _maybe_ return an error channel, to notify if the server fails somehow
func (c *Sys) Monitor(loglevel string, logJSON bool, stopCh chan struct{}, q bool) (chan string, error) {
	r := c.c.NewRequest("GET", "/v1/sys/monitor")
	// TODO: I need to set the query options here

	ctx, cancelFunc := context.WithCancel(context.Background())
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}

	logCh := make(chan string, 64)
	go func() {
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)

		for {
			select {
			case <-stopCh:
				close(logCh)
				cancelFunc()
				return
			case <-ctx.Done():
				stopCh <- struct{}{}
				close(logCh)
				cancelFunc()
				return
			default:
			}

			if scanner.Scan() {
				// An empty string signals to the caller that
				// the scan is done, so make sure we only emit
				// that when the scanner says it's done, not if
				// we happen to ingest an empty line.
				if text := scanner.Text(); text != "" {
					logCh <- text
				} else {
					logCh <- " "
				}
			} else {
				// If Scan() returns false, that means the context deadline was exceeded, so
				// terminate this routine and start a new request.
				stopCh <- struct{}{}
				close(logCh)
				cancelFunc()
				return
			}
		}
	}()

	return logCh, nil
}
