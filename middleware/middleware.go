package middleware

import (
	"github.com/valyala/fasthttp"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"
)

// Logging is middleware for logging information about each request.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, req)
		log.Printf("%s %s %s", req.Method, req.RequestURI, time.Since(start))
	})
}

// PanicRecovery is middleware for recovering from panics in `next` and
// returning a StatusInternalServerError to the client.
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				log.Println(string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, req)
	})
}

// ---------------------------------------------------------------
// For fasthttp modify from https://github.com/AubSs/fasthttplogger

var (
	output = log.New(os.Stdout, "", 0)
)

var (
	green  = string([]byte{27, 91, 48, 48, 58, 51, 50, 109})
	yellow = string([]byte{27, 91, 48, 48, 59, 51, 51, 109})
	red    = string([]byte{27, 91, 48, 48, 59, 51, 49, 109})
	blue   = string([]byte{27, 91, 48, 48, 59, 51, 52, 109})
	white  = string([]byte{27, 91, 48, 109})
)

func getColorByStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return green
	case code >= 300 && code < 400:
		return blue
	case code >= 400 && code < 500:
		return yellow
	default:
		return red
	}
}

func colorStatus(code int) string {
	return getColorByStatus(code) + strconv.Itoa(code) + white
}

func colorMethod(method []byte, code int) string {
	return getColorByStatus(code) + string(method) + white
}

func getHttp(ctx *fasthttp.RequestCtx) string {
	if ctx.Response.Header.IsHTTP11() {
		return "HTTP/1.1"
	}
	return "HTTP/1.0"
}

// LoggerAndPanicRecover middleware for logging and panic recover
func LoggerAndPanicRecover(req fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Println(string(debug.Stack()))
				ctx.Response.SetStatusCode(fasthttp.StatusInternalServerError)
				return
			}
		}()

		begin := time.Now()
		req(ctx)
		end := time.Now()
		output.Printf("[%v] %v | %s | %s %s - %v - %v | %s",
			end.Format("2006/01/02 - 15:04:05"),
			ctx.RemoteAddr(),
			getHttp(ctx),
			colorMethod(ctx.Method(), ctx.Response.Header.StatusCode()),
			ctx.RequestURI(),
			colorStatus(ctx.Response.Header.StatusCode()),
			end.Sub(begin),
			ctx.UserAgent(),
		)
	})
}
