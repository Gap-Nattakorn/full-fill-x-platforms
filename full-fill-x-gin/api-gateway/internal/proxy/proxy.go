package proxy

import (
	"net/http/httputil"
	"net/url"

	"github.com/GapNattakorn/full-fill-x-platforms/full-fill-x-gin/api-gateway/internal/response"
	"github.com/gin-gonic/gin"
)

func Forward(target string) gin.HandlerFunc {
	targetURL, err := url.Parse(target)

	return func(c *gin.Context) {
		if err != nil {
			response.Error(c, 500, "PROXY_CONFIG_ERROR", "Invalid upstream service URL")
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(targetURL)

		c.Request.URL.Scheme = targetURL.Scheme
		c.Request.URL.Host = targetURL.Host
		c.Request.Host = targetURL.Host
		c.Request.Header.Set("X-Request-ID", c.GetString("request_id"))

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}