package middleware

import (
	"github.com/Beretta350/authentication/internal/pkg/dto"
	"github.com/Beretta350/authentication/pkg/csrf"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func CSRFHandler(wrapper csrf.CSRFWrapper) gin.HandlerFunc {
	return func(c *gin.Context) {
		if wrapper.IsIgnoredPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		session := sessions.Default(c)
		c.Set(csrf.Secret, wrapper.GetSecret())

		s, ok := session.Get(csrf.Session).(string)
		if !ok || len(s) == 0 {
			defaultCSRFErrorFunc(c)
			return
		}

		token := defaultGetToken(c)

		valid := wrapper.ValidateToken(token, s)
		if !valid {
			defaultCSRFErrorFunc(c)
			return
		}
	}
}

func defaultCSRFErrorFunc(c *gin.Context) {
	response := dto.ForbiddenResponse("Invalid CSRF token", nil)
	c.JSON(response.StatusCode, response)
	c.Abort()
}

func defaultGetToken(c *gin.Context) string {
	r := c.Request

	if t := r.FormValue("_csrf"); len(t) > 0 {
		return t
	} else if t := r.URL.Query().Get("_csrf"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-CSRF-TOKEN"); len(t) > 0 {
		return t
	} else if t := r.Header.Get("X-XSRF-TOKEN"); len(t) > 0 {
		return t
	}

	return ""
}
