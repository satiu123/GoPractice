package geecache

import (
	"fmt"
	"log"
	"net/http"

	gee "gee-web"
)

const defaultBasePath = "/_geecache/"

// HTTPPool 为 HTTP 对等点池实现了 PeerPicker。
type HTTPPool struct {
	// this peer's base URL, e.g. "https://example.net:8000"
	self     string
	basePath string
}

// NewHTTPPool 初始化一个 HTTP 对等点池。
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log 使用服务器名称记录信息
func (p *HTTPPool) Log(format string, v ...any) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// RegisterRoutes 将 geecache 的路由注册到 gee-web 引擎
func (p *HTTPPool) RegisterRoutes(engine *gee.Engine) {
	cacheGroup := engine.Group(p.basePath)
	{
		cacheGroup.GET("/:groupname/:key", p.get)
	}
}

func (p *HTTPPool) get(c *gee.Context) {
	groupName := c.Param("groupname")
	key := c.Param("key")

	group := GetGroup(groupName)
	if group == nil {
		c.String(http.StatusNotFound, "no such group: "+groupName)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.SetHeader("Content-Type", "application/octet-stream")
	c.Data(http.StatusOK, view.ByteSlice())
}
