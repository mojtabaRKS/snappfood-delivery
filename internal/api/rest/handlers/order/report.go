package order

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *OrderHandler) ReportDelay(c *gin.Context) {
	var req struct {
		OrderId uint `json:"order_id"`
	}

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "order_id is invalid",
		})
		return
	}

	isReportable, err := h.orderService.IsReportable(c, req.OrderId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "order not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if isReportable {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "order delivery time is not over yet",
		})
		return
	}

	isProcessing, err := h.orderService.IsProcessing(c, req.OrderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if isProcessing {
		c.JSON(http.StatusConflict, gin.H{
			"message": "order is already reported",
		})
		return
	}

	deliveryTime, err := h.orderService.ReportDelay(c, req.OrderId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("an error occured %s", err.Error()),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data": gin.H{
			"deliveery_at": deliveryTime,
		},
	})
}
