package service

import (
	"encoding/json"
	"fmt"
	"github.com/ExchangeUnion/xud-docker-api-poc/utils"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hpcloud/tail"
	"io"
	"net/http"
	"strings"
)

func (t *Manager) ConfigureRouter(r *gin.Engine) {
	r.Use(static.Serve("/", static.LocalFile("/ui", false)))

	api := r.Group("/api")
	{
		api.GET("/v1/services", func(c *gin.Context) {
			var result []ServiceEntry

			result = append(result, ServiceEntry{"xud", "XUD"})
			result = append(result, ServiceEntry{"lndbtc", "LND (Bitcoin)"})
			result = append(result, ServiceEntry{"lndltc", "LND (Litecoin)"})
			result = append(result, ServiceEntry{"connext", "Connext"})
			result = append(result, ServiceEntry{"bitcoind", "Bitcoind"})
			result = append(result, ServiceEntry{"litecoind", "Litecoind"})
			result = append(result, ServiceEntry{"geth", "Geth"})
			result = append(result, ServiceEntry{"arby", "Arby"})
			result = append(result, ServiceEntry{"boltz", "Boltz"})
			result = append(result, ServiceEntry{"webui", "Web UI"})

			c.JSON(http.StatusOK, result)
		})

		api.GET("/v1/status", func(c *gin.Context) {
			status := t.GetStatus()

			var result []ServiceStatus

			for _, svc := range t.services {
				result = append(result, ServiceStatus{Service: svc.GetName(), Status: status[svc.GetName()]})
			}

			c.JSON(http.StatusOK, result)
		})

		api.GET("/v1/status/:service", func(c *gin.Context) {
			service := c.Param("service")
			s, err := t.GetService(service)
			if err != nil {
				utils.JsonError(c, err.Error(), http.StatusNotFound)
				return
			}
			status, err := s.GetStatus()
			if err != nil {
				status = fmt.Sprintf("Error: %s", err)
			}
			c.JSON(http.StatusOK, ServiceStatus{Service: service, Status: status})
		})

		api.GET("/v1/logs/:service", func(c *gin.Context) {
			service := c.Param("service")
			s, err := t.GetService(service)
			if err != nil {
				utils.JsonError(c, err.Error(), http.StatusNotFound)
				return
			}
			since := c.DefaultQuery("since", "1h")
			tail := c.DefaultQuery("tail", "all")
			logs, err := s.GetLogs(since, tail)
			if err != nil {
				utils.JsonError(c, err.Error(), http.StatusInternalServerError)
				return
			}
			c.Header("Content-Type", "text/plain")
			c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.log\"", service))
			for _, line := range logs {
				_, err = c.Writer.WriteString(line + "\n")
				if err != nil {
					utils.JsonError(c, err.Error(), http.StatusInternalServerError)
				}
			}
		})

		api.GET("/v1/setup-status", func(c *gin.Context) {
			c.Stream(func(w io.Writer) bool {
				logfile := fmt.Sprintf("/root/network/logs/%s.log", t.network)
				t, err := tail.TailFile(logfile, tail.Config{
					Follow: true,
					ReOpen: true})
				if err != nil {
					return false
				}
				for line := range t.Lines {
					if strings.Contains(line.Text, "Waiting for XUD dependencies to be ready") {
						status := SetupStatus{Status: "Waiting for XUD dependencies to be ready", Details: nil}
						j, _ := json.Marshal(status)
						c.Writer.Write(j)
						c.Writer.Write([]byte("\n"))
						c.Writer.Flush()
					} else if strings.Contains(line.Text, "LightSync") {
						parts := strings.Split(line.Text, " [LightSync] ")
						parts = strings.Split(parts[1], " | ")
						details := map[string]string{}
						status := SetupStatus{Status: "Syncing light clients", Details: details}
						for _, p := range parts {
							kv := strings.Split(p, ": ")
							details[kv[0]] = kv[1]
						}
						j, _ := json.Marshal(status)
						c.Writer.Write(j)
						c.Writer.Write([]byte("\n"))
						c.Writer.Flush()
					} else if strings.Contains(line.Text, "Setup wallets") {
						status := SetupStatus{Status: "Setup wallets", Details: nil}
						j, _ := json.Marshal(status)
						c.Writer.Write(j)
						c.Writer.Write([]byte("\n"))
						c.Writer.Flush()
					} else if strings.Contains(line.Text, "Create wallets") {
						status := SetupStatus{Status: "Create wallets", Details: nil}
						j, _ := json.Marshal(status)
						c.Writer.Write(j)
						c.Writer.Write([]byte("\n"))
						c.Writer.Flush()
					} else if strings.Contains(line.Text, "Restore wallets") {
						status := SetupStatus{Status: "Restore wallets", Details: nil}
						j, _ := json.Marshal(status)
						c.Writer.Write(j)
						c.Writer.Write([]byte("\n"))
						c.Writer.Flush()
					} else if strings.Contains(line.Text, "Setup backup location") {
						status := SetupStatus{Status: "Setup backup location", Details: nil}
						j, _ := json.Marshal(status)
						c.Writer.Write(j)
						c.Writer.Write([]byte("\n"))
						c.Writer.Flush()
					} else if strings.Contains(line.Text, "Unlock wallets") {
						status := SetupStatus{Status: "Unlock wallets", Details: nil}
						j, _ := json.Marshal(status)
						c.Writer.Write(j)
						c.Writer.Write([]byte("\n"))
						c.Writer.Flush()
					} else if strings.Contains(line.Text, "Start shell") {
						break
					}
				}
				return false
			})
		})
	}

	for _, svc := range t.services {
		svc.ConfigureRouter(api)
	}
}
