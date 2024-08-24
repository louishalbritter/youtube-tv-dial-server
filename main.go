package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	uuid := uuid.New().String()
	ip := getLocalIP()

	t := &Template{
		templates: template.Must(template.ParseGlob("./templates/**")),
	}

	friendlyName, err := os.Hostname()
	if err != nil {
		friendlyName = ip
	}

	service := Service{
		ModelName:    "YouTube",
		Uuid:         uuid,
		Manufacture:  "YouTube",
		FriendlyName: friendlyName,
		BaseUrl:      fmt.Sprintf("http://%v:3000/dial", ip),
	}

	go service.ssdp()

	e := echo.New()

	e.Use(middleware.Logger())

	e.HideBanner = true

	e.Renderer = t

	e.GET("/dial/ssdp/device-desc.xml", func(c echo.Context) error {
		c.Response().Header().Set("Application-URL", service.BaseUrl+"/apps")
		c.Response().Header().Set("Content-Type", "application/xml")

		return c.Render(http.StatusOK, "device-desc.xml", service)
	})

	e.GET("/dial/apps/YouTube", func(c echo.Context) error {
		c.Response().Header().Set("Content-Type", "application/xml")
		c.Response().Header().Set("LOCATION", "http://"+ip+":3000/dial/apps/YouTube/test")

		return c.Render(http.StatusOK, "application.xml", service)
	})

	e.POST("/dial/apps/YouTube", func(c echo.Context) error {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			log.Fatal(err)
		}

		go service.start(string(bodyBytes))

		c.Response().Header().Set("LOCATION", "http://"+ip+":3000/dial/apps/YouTube/test")
		return c.NoContent(http.StatusOK)
	})

	e.DELETE("/dial/apps/YouTube", func(c echo.Context) error {
		service.stop()

		return c.NoContent(http.StatusOK)
	})

	e.Logger.Fatal(e.Start(":3000"))
}
