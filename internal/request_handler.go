package internal

import (
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func echoHandlerWith(verbose bool, verboseHeaders bool, delay int32) func(*gin.Context) {
	maybeDelay := func() {
		if delay > 0 {
			if verbose {
				log.Print(Green("Delaying response in %d millis", delay))
			}
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	return handler(verbose, verboseHeaders, maybeDelay)
}

func handler(verbose bool, verboseHeaders bool, doBefore func()) func(*gin.Context) {
	return func(c *gin.Context) {
		// todo: in case verbose is off, we definitely don't have to read the entire body into memory.
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(Red("Failed to read request body: %s", err.Error()))
		} else {
			if verboseHeaders {
				printHeaders(c)
			}

			if verbose && len(bodyBytes) > 0 {
				printBody(bodyBytes)
			}
		}

		doBefore()

		requestContentType := c.GetHeader("Content-Type")
		if requestContentType != "" {
			c.Writer.Header().Set("Content-Type", c.GetHeader("Content-Type"))
		}
		_, err = c.Writer.Write(bodyBytes)
		if err != nil {
			log.Println(Red("Failed to echo request body: %s", err.Error()))
		}
	}
}

func printHeaders(ctx *gin.Context) {
	var headersStr = ""
	for header := range ctx.Request.Header {
		headersStr += fmt.Sprintf(
			"%s : %s\n\r", Bold("%s", header), Bold("%s", Cyan("%s", ctx.Request.Header.Get(header))))
	}
	log.Printf(`Received headers:

%s`, headersStr)
}

func printBody(bodyBytes []byte) {
	// todo: this ignore encoding and assumes the body is printable.. maybe we can do better.
	log.Printf(`Received body:

%s

`,
		Yellow("%s", strings.TrimSpace(strings.ReplaceAll(string(bodyBytes), "%", "%%"))),
	)
}
