package vendor

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *VendorHandler) GetDelayReport(c *gin.Context) {
	vendors, err := h.orderService.GetDelayReport(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"data":    vendors,
	})
}
