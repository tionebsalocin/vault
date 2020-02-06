package http

import (
	"fmt"
	"io"
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	log "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/vault/command/monitor"
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

		isJson := core.SanitizedConfig()["log_format"] == "json"
		monitor := monitor.New(512, core.Logger(), &log.LoggerOptions{
			Level:      logLevel,
			JSONFormat: isJson,
		})

		w.Header().Set("Content-Type", "application/json")

		logCh := monitor.Start()
		defer monitor.Stop()
		errCh := make(chan error, 2)

		go func() {
			for {
				select {
				case log := <-logCh:
					tmp := string(log[:])

					// We got text back because JSON logging isn't enabled, but we'd like to return things from this function in a consistent manner.
					if !isJson {
						r, _ := regexp.Compile("^([0-9\.:T-])+ \[([A-Z]+)\] ([a-z\.-]): (.+)$")
						matches := r.FindAllStringSubmatch(tmp, -1)
						fmt.Printf("matches = %v\n", matches)
						output := map[string]interface{} {
							"@level": matches[1],
							"@message": matches[3],
							"@module": matches[2],
							"@timestamp": matches[0],
						}

						json, err := json.Marshal(output)
						if err != nil {
							respondError(w, http.StatusInternalServerError, err)
						}

						tmp = string(json)
					}

					if _, err := w.Write(tmp); err != nil {
						errCh <- err
						return
					}

					if f, ok := w.(http.Flusher); ok {
						f.Flush()
					}
				}
			}
		}()

		e := <-errCh

		// this feels mostly useless?
		if e != nil &&
			(e == io.EOF ||
				strings.Contains(e.Error(), "closed") ||
				strings.Contains(e.Error(), "EOF")) {
			fmt.Println("client closed the connection")
			e = nil
		}
	})
}

