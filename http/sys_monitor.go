package http

import (
	"fmt"
	"github.com/hashicorp/vault/command/monitor"
	"net/http"
	"time"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/vault"
)

func handleSysMonitor(core *vault.Core) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ll := r.URL.Query().Get("log_level")
		if ll == "" {
			ll = "INFO"
		}
		logLevel := log.LevelFromString(ll)

		if logLevel == log.NoLevel {
			respondError(w, http.StatusBadRequest, fmt.Errorf("invalid log level"))
			return
		}

		flusher, ok := w.(http.Flusher)
		if !ok {
			respondError(w, http.StatusBadRequest, fmt.Errorf("streaming not supported"))
			return
		}

		isJson := core.SanitizedConfig()["log_format"] == "json"
		monitor := monitor.New(512, core.Logger(), &log.LoggerOptions{
			Level:      logLevel,
			JSONFormat: isJson,
		})

		//w.Header().Set("Content-Type", "application/json")
		//
		logCh := monitor.Start()
		//defer monitor.Stop()
		//errCh := make(chan error, 2)

		w.WriteHeader(http.StatusOK)

		// 0 byte write is needed before the Flush call so that if we are using
		// a gzip stream it will go ahead and write out the HTTP response header
		w.Write([]byte(""))
		flusher.Flush()

		// Stream logs until the connection is closed.
		for {
			select {
			case <-r.Context().Done():
				monitor.Stop()
				return
			case log := <-logCh:
				fmt.Fprint(w, string(log))
				flusher.Flush()
			}
		}

		//go func() {
		//	for {
		//		select {
		//		case log := <-logCh:
		//			tmp := log
		//
		//			// We got text back because JSON logging isn't enabled, but we'd like to return things from this function in a consistent manner.
		//			// So, parse it all out, and form it into JSON
		//			if !isJson {
		//				r, _ := regexp.Compile(`^([0-9:T\.-]+)\s+\[([A-Z]+)\]\s+([a-z\.-]+):\s+(.+)$`)
		//				matches := r.FindAllStringSubmatch(string(tmp), -1)
		//				fmt.Printf("matches = %v\n", matches)
		//				output := map[string]interface{}{
		//					"@level":     matches[1],
		//					"@message":   matches[3],
		//					"@module":    matches[2],
		//					"@timestamp": matches[0],
		//				}
		//
		//				json, err := json.Marshal(output)
		//				if err != nil {
		//					respondError(w, http.StatusInternalServerError, err)
		//				}
		//				tmp = json
		//			}
		//
		//			if _, err := w.Write(tmp); err != nil {
		//				errCh <- err
		//				return
		//			}
		//
		//			if f, ok := w.(http.Flusher); ok {
		//				f.Flush()
		//			}
		//		}
		//	}
		//}()

		//e := <-errCh

		// this is useless?
		//if e != nil &&
		//	(e == io.EOF ||
		//		strings.Contains(e.Error(), "closed") ||
		//		strings.Contains(e.Error(), "EOF")) {
		//	fmt.Println("client closed the connection")
		//	e = nil
		//}
	})
}

type MonitorResponse struct {
	Level     log.Level `json:"@level"`
	Message   string    `json:"@message"`
	Module    string    `json:"@module"`
	Timestamp time.Time `json:"@timestamp"`
}
