package gemini

import "github.com/gin-gonic/gin"

func Url(r *gin.Engine) {
	v1 := r.Group("/v1/gemini")
	{
		v1.POST("", GetGeminiInteraction)
	}
}
