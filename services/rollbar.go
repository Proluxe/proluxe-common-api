package services

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
	"github.com/rollbar/rollbar-go"
	u "github.com/scottraio/go-utils"
)

func ErrorHandling() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rval := recover(); rval != nil {
				debug.PrintStack()

				// Wrap rval into error if not already one
				var err error
				switch e := rval.(type) {
				case error:
					err = e
				default:
					err = fmt.Errorf("panic: %v", e)
				}

				rollbar.SetToken(u.GetDotEnvVariable("ROLLBAR_API_KEY"))
				rollbar.SetEnvironment("development")

				rollbar.RequestErrorWithStackSkipWithExtras(rollbar.CRIT, c.Request, err, 2, map[string]interface{}{
					"endpoint": c.Request.RequestURI,
				})

				c.AbortWithStatus(http.StatusInternalServerError)
			}
		}()

		c.Next()
	}
}
