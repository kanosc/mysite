package router

import (
	"fmt"

	"net/http"

	"github.com/kanosc/mysite/server/endpoints"
	"github.com/kanosc/mysite/server/service"

	"github.com/kanosc/mysite/server/transport"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
)

var r = gin.Default()

func PrintRouter() {
	fmt.Println("This is router")
}

func DefaultRoute(c *gin.Context) {
	c.HTML(200, "index.html", nil)
}

func initStaticServer() {
	r.StaticFS("/assets/", http.Dir("dist/assets"))
	r.StaticFS("/images/", http.Dir("dist/images"))
	r.StaticFS("/.well-known", http.Dir(".well-known"))
	r.LoadHTMLFiles("./dist/index.html", "upload.html")
	r.GET("/", DefaultRoute)
	r.NoRoute(DefaultRoute)
	r.GET("/upload", func(c *gin.Context) { c.HTML(200, "upload.html", nil) })
}

func initMicroServiceRouter() {
	var flsvc service.FileAccessor
	flsvc = new(service.CommonFileAccessor)
	flvcHandler := httptransport.NewServer(
		endpoints.MakeFileListEndpoint(flsvc),
		transport.DecodeFileListRequest,
		transport.EncodeFileListResponse,
	)
	// test commond: curl -XPOST -d'{"DirName":"testdir"}' localhost:5173/files
	fileRoute := r.Group("/api/file/v1/")
	{
		fileRoute.GET("/files", gin.WrapH(flvcHandler))
		fileRoute.GET("/:dirname/:filename", service.FileRequestChecker(), service.DownloadFile)
		fileRoute.DELETE("/:dirname/:filename", service.FileRequestChecker(), service.DeleteFile)
		fileRoute.POST("/:dirname", service.UploadFile)
	}

	chatRoute := r.Group("/api/chat/v1/")
	{
		chatRoute.GET("/rooms", service.QueryRooms)
		chatRoute.GET("/wschat/:roomname", service.HandlePostMsg)
		chatRoute.POST("/room", service.CreateRoom)
		chatRoute.DELETE("/:roomname", service.DeleteRoom)
		chatRoute.GET("/msg/:roomname", service.PasswordChecker(), service.GetRoomMessages)
		chatRoute.POST("/auth/:roomname", service.AuthToken)
	}

	/*
	kmsRoute := r.Group("/api/kms/v1/")
	{
		kmsRoute.GET("/:keyid", service.KeyOperationAuth(), service.KeyRequestChecker(), service.HandleGetClientKeyById)
		kmsRoute.POST("/:clientid/key", service.KeyCreateRateLimiter(20), service.KeyOperationAuth(), service.HandleCreateClientKey)
	}
	*/

}

func InitRouter(mode string) {
	service.RedisInit(mode)
	initStaticServer()
	initMicroServiceRouter()
}

func Start(addr string) {
	r.Run(addr)
}
