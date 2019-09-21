package httpserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/herb-go/util"
)

//RecoverMiddleware create recover middleware with given logger.
func RecoverMiddleware(logger *log.Logger) func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if logger == nil {
		logger = log.New(os.Stderr, log.Prefix(), log.Flags())
	}
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		defer func() {
			if r := recover(); r != nil {
				err := r.(error)
				var result string
				if util.IsErrorIgnored(err) == false {
					lines := strings.Split(string(debug.Stack()), "\n")
					length := len(lines)
					maxLength := util.LoggerMaxLength*2 + 7
					if length > maxLength {
						length = maxLength
					}
					var output = make([]string, length-6)
					output[0] = fmt.Sprintf("Panic: %s - http request %s \"%s\" ", err.Error(), req.Method, req.URL.String())
					output[0] += "\n" + lines[0]
					copy(output[1:], lines[7:])
					result = strings.Join(output, "\n")
					logger.Println(result)
				}
				if util.Debug {
					http.Error(w, result, http.StatusInternalServerError)
				} else {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}
		}()
		next(w, req)
	}
}
