package gin

import (
	"github.com/ammorteza/clean_architecture/http"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ginRouter struct {
	ginDispatcher 		*gin.Engine
}

func New() router.Router {
	return &ginRouter{
		ginDispatcher: gin.Default(),
	}
}

func (gr *ginRouter)GetDispatcher() interface{}{
	return gr.ginDispatcher
}

func (gr *ginRouter)GET(uri string, f func(w http.ResponseWriter, r *http.Request)){
	gr.ginDispatcher.GET(uri, func(context *gin.Context) {
		temp := context.Request.URL.Query()
		for _, param := range context.Params{
			temp.Add(param.Key, param.Value)
		}
		context.Request.URL.RawQuery = temp.Encode()
		f(context.Writer, context.Request)
	})
}

func (gr *ginRouter)POST(uri string, f func(w http.ResponseWriter, r *http.Request)){
	gr.ginDispatcher.POST(uri, gin.WrapF(f))
}

func (gr *ginRouter)SERVE(port string){
	gr.ginDispatcher.Run(":8080")
}