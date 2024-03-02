package order

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *OrderHandler) Proccess(c *gin.Context) {
	str := c.Request.Header.Get("employee_id")
	if str == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "invalid employee_id in header",
		})
		return
	}

	employeeId, err := strconv.Atoi(str)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": fmt.Sprintf("cant convert to int. err : %s", err.Error()),
		})
		return
	}

	isProcessing, err := h.orderService.IsProccessing(c, employeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	if isProcessing {
		c.JSON(http.StatusConflict, gin.H{
			"message": "cant start new order withount finish previuos",
		})
		return
	}

	order, err := h.orderService.Proccess(c, employeeId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    order,
	})
}
