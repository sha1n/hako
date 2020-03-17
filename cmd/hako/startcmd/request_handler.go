package startcmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sha1n/hako/cmd/hako/console"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

func EchoHandlerWith(verbose bool, verboseHeaders bool, delay int32) func(*gin.Context) {
	maybeDelay := func() {
		if delay > 0 {
			log.Printf(console.Green("Delaying response in %d millis", delay))
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}
	}

	return handler(verbose, verboseHeaders, maybeDelay)
}

func handler(verbose bool, verboseHeaders bool, doBefore func()) func(*gin.Context) {
	return func(c *gin.Context) {
		// todo: in case verbose is off, we definitely don't have to read the entire body into memory.
		bodyBytes, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(console.Red("Failed to read request body: %s", err.Error()))
		} else {
			if verbose && len(bodyBytes) > 0 {
				printBody(bodyBytes)
			}

			if verboseHeaders {
				printHeaders(c)
			}
		}

		doBefore()

		_, err = c.Writer.Write(bodyBytes)
		if err != nil {
			log.Println(console.Red("Failed to echo request body: %s", err.Error()))
		}
	}
}

func printHeaders(ctx *gin.Context) {
	var headersStr = ""
	for header := range ctx.Request.Header {
		headersStr += fmt.Sprintf(
			"%s : %s\n\r", console.Bold(header), console.Bold(console.Cyan(ctx.Request.Header.Get(header))))
	}
	log.Printf(`Received headers:

%s`, headersStr)
}

func printBody(bodyBytes []byte) {
	// todo: this ignore encoding and assumes the body is printable.. maybe we can do better.
	log.Printf(`Received body:

%s

`,
		console.Yellow(strings.TrimSpace(string(bodyBytes))),
	)
}
