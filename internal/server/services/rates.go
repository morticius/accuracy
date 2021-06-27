package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/morticius/accuracy/pkg/services"
	"log"
	"net/http"
)

type RatesService struct {
	btcService *services.BTCService
}

func NewRatesService(b *services.BTCService) *RatesService {
	return &RatesService{
		btcService: b,
	}
}

func (r *RatesService) GetRateBTCToUAHHandler(c *gin.Context) {
	rate, err := r.btcService.GetPriceBTCInUAH()
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}

	c.JSON(http.StatusOK, gin.H{
		"rate": fmt.Sprintf("%.2f", rate),
	})
}