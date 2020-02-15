package api

import (
	"bufio"
	"fmt"
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
func (c *Sys) Monitor(loglevel string, logJSON bool, stopCh <-chan struct{}, q bool) (chan string, error) {
	r := c.c.NewRequest("GET", "/v1/sys/monitor")
	// TODO: I need to set the query options here

	ctx, cancelFunc := context.WithCancel(context.Background())
	//defer cancelFunc()
	resp, err := c.c.RawRequestWithContext(ctx, r)
	if err != nil {
		return nil, err
	}

	//var result LeaderResponse
	//err = resp.DecodeJSON(&result)
	//return &result, err


	// ===================================

	logCh := make(chan string, 64)
	go func() {
		defer resp.Body.Close()
		scanner := bufio.NewScanner(resp.Body)
		//reader := bufio.NewReader(resp.Body)

		for {
			select {
			case <-stopCh:
				fmt.Println("input on the stopCh. stopping")
				close(logCh)
				cancelFunc()
				return
			default:
			}

			//line, err := reader.ReadBytes('\n')
			//
			//if err != nil {
			//	fmt.Printf("got an error: %v\n", err)
			//	close(logCh)
			//	cancelFunc()
			//	return
			//}

			//logCh <- string(line)

			if scanner.Scan() {
				fmt.Println("scanned something")
				// An empty string signals to the caller that
				// the scan is done, so make sure we only emit
				// that when the scanner says it's done, not if
				// we happen to ingest an empty line.
				if text := scanner.Text(); text != "" {
					fmt.Printf("text = |%v|\n", text)
					fmt.Println("sending this over logCh")
					logCh <- text
				} else {
					fmt.Println("sending a space over logCh")
					logCh <- " "
				}
			} else {
				fmt.Println("Scan() returned false. sending an empty string over logCh")
				logCh <- ""
			}
		}
	}()

	return logCh, nil
}

// Monitor returns a channel which will receive streaming logs. Providing a
// non-nil stopCh can be used to close the connection and stop streaming.
//func (c *Sys) Monitor(stopCh <-chan struct{}) (<-chan *MonitorResponse, <-chan error) {
//	errCh := make(chan error, 1)
//
//	// TODO: I need to somehow pass query options to this URL
//	r := c.c.NewRequest("GET", "/v1/sys/monitor")
//	lines := make(chan *MonitorResponse, 10)
//
//	fmt.Println("awesome fun times")
//
//	go func() {
//		f, _ := os.Create("good-times.txt")
//		defer f.Close()
//
//		ctx, cancelFunc := context.WithCancel(context.Background())
//		defer cancelFunc()
//		resp, err := c.c.RawRequestWithContext(ctx, r)
//
//		if err != nil {
//			errCh <- err
//			return
//		}
//
//		f.WriteString("no errors\n")
//
//		defer resp.Body.Close()
//		dec := json.NewDecoder(resp.Body)
//
//		f.WriteString("awesome\n")
//
//		buf := new(bytes.Buffer)
//		buf.ReadFrom(resp.Body)
//		bodyString := buf.String()
//
//		//bodyBytes, _ := ioutil.ReadAll(resp.Body)
//		f.WriteString("dope\n")
//		//bodyString := string(bodyBytes)
//
//		f.WriteString(fmt.Sprintf("body = |%v|\n", bodyString))
//
//		for {
//			select {
//			case <- stopCh:
//				f.WriteString("received something on the stop channel\n")
//				close(lines)
//				return
//			default:
//			}
//
//			var line MonitorResponse
//			if err := dec.Decode(&line); err != nil {
//				f.WriteString(fmt.Sprintf("error decoding the body: %v\n", err))
//				close(lines)
//				errCh <- err
//				return
//			}
//
//			f.WriteString(fmt.Sprintf("line = %v\n", line))
//
//			lines <- &line
//		}
//	}()
//
//	return lines, errCh
//}
//
//type MonitorResponse struct {
//	Level     log.Level `json:"@level"`
//	Message   string    `json:"@message"`
//	Module    string    `json:"@module"`
//	Timestamp time.Time `json:"@timestamp"`
//}
//
//func (m MonitorResponse) String() string {
//	return fmt.Sprint("%v [%v] %v: %v", m.Timestamp, m.Level, m.Module, m.Message)
//}
