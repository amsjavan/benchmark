package fast

import (
	"fmt"
	"github.com/valyala/fasthttp"
	"time"
)

type MyHandler struct {
	foobar string
}


// request handler in fasthttp style, i.e. just plain function.
func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	time.Sleep(200 * time.Millisecond)
	_, _ = fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
}

func init(){

}


func RunFastServer(){
	// pass plain function to fasthttp
	_ = fasthttp.ListenAndServe(":8080", fastHTTPHandler)


}