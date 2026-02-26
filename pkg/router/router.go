package router

import (
	"app/pkg/controller"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	//r.StaticFS("/", http.Dir(filepath.Join(dir, "./html")))

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.GET("/cookie", controller.HandleCookie)

	r.GET("/", func(c *gin.Context) {
		c.File("./html/index.html")
	})

	r.Static("/static", "./static")

	err := RegisterHTMLRoutes(r, "./html", "/")
	if err != nil {
		panic("Failed to register HTML routes: " + err.Error())
	}
}

func RegisterHTMLRoutes(router *gin.Engine, baseDir string, prefix string) error {
	files, err := os.ReadDir(baseDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		path := filepath.Join(baseDir, file.Name())
		if file.IsDir() {
			err := RegisterHTMLRoutes(router, path, prefix+file.Name()+"/")
			if err != nil {
				return err
			}
		} else if strings.HasSuffix(strings.ToLower(file.Name()), ".html") {
			urlPath := prefix + strings.TrimSuffix(file.Name(), ".html")
			router.GET(urlPath, func(c *gin.Context) {
				c.File(path)
			})
		}
	}
	return nil
}
