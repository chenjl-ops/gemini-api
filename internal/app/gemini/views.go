package gemini

import (
	"fmt"
	"gemini-api/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGeminiInteraction(c *gin.Context) {
	var req AiRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.WithError(c, http.StatusBadRequest, fmt.Sprintf("invalid input %v", err))
	} else {
		results, err := CallGeminiInteraction(req.ApiKey, req.Prompt)
		if err != nil {
			response.WithError(c, http.StatusInternalServerError, err.Error())
		}
		response.WithResult(c, http.StatusOK, "success", results)
	}

}
