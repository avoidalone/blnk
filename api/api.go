package api

import (
	"fmt"

	"github.com/jerry-enebeli/blnk/config"

	"github.com/jerry-enebeli/blnk/api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/jerry-enebeli/blnk"
)

type Api struct {
	blnk   *blnk.Blnk
	router *gin.Engine
}

func (a Api) Router() *gin.Engine {
	router := a.router
	router.POST("/ledgers", a.CreateLedger)
	router.GET("/ledgers/:id", a.GetLedger)
	router.GET("/ledgers", a.GetAllLedgers)

	router.POST("/balances", a.CreateBalance)
	router.GET("/balances/:id", a.GetBalance)

	router.POST("/balance-monitors", a.CreateBalanceMonitor)
	router.GET("/balance-monitors/:id", a.GetBalanceMonitor)
	router.GET("/balance-monitors", a.GetAllBalanceMonitors)
	router.GET("/balance-monitors/balances/:balance_id", a.GetBalanceMonitorsByBalanceID)
	router.PUT("/balance-monitors/:id", a.UpdateBalanceMonitor)

	router.POST("/transactions", a.QueueTransaction)
	router.POST("/refund-transaction/:id", a.RefundTransaction)
	router.GET("/transactions/:id", a.GetTransaction)
	router.PUT("/transactions/inflight/:id/:status", a.UpdateInflightStatus)

	router.POST("/identities", a.CreateIdentity)
	router.GET("/identities/:id", a.GetIdentity)
	router.PUT("/identities/:id", a.UpdateIdentity)
	router.GET("/identities", a.GetAllIdentities)

	router.POST("/accounts", a.CreateAccount)
	router.GET("/accounts/:id", a.GetAccount)
	router.GET("/accounts", a.GetAllAccounts)

	router.GET("/mocked-account", a.generateMockAccount)

	router.GET("/backup", a.BackupDB)
	router.GET("/backup-s3", a.BackupDBS3)
	return a.router
}

func NewAPI(b *blnk.Blnk) *Api {
	gin.SetMode(gin.ReleaseMode)
	conf, err := config.Fetch()
	if err != nil {
		return nil
	}
	r := gin.Default()
	if conf.Server.Secure {
		r.Use(middleware.SecretKeyAuthMiddleware())
	}

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, "server running...")
	})

	r.POST("/webhook", func(c *gin.Context) {
		var payload map[string]interface{}
		err := c.Bind(&payload)
		if err != nil {
			return
		}
		fmt.Println(payload)
	})

	return &Api{blnk: b, router: r}
}
