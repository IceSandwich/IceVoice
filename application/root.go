package application

import (
	"context"
	"fmt"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
	"github.com/lxn/walk"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	ApiVersion = "v1.0.0"
)

type RouteMethod int

const (
	RouteMethodGet RouteMethod = iota
	RouteMethodPost
)

type notifyItem struct {
	Name   string
	Action func()
}

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	Port int `help:"Port to listen on" short:"p" default:"8888"`
}

type Application struct {
	window        *walk.MainWindow
	notify        *walk.NotifyIcon
	notifyActions []notifyItem

	router *chi.Mux
	api    huma.API
	port   int

	isInit bool
}

var app Application

func initapp() {
	if app.isInit {
		return
	}

	// Create a new router & API.
	app.router = chi.NewMux()
	app.api = humachi.New(app.router, huma.DefaultConfig("IceVoice", ApiVersion))
	app.notifyActions = make([]notifyItem, 0)
	app.isInit = true
}

func init() {
	initapp()
}

func RegisterGetRouter[I, O any](path string, handler func(context.Context, *I) (*O, error)) {
	if app.isInit == false {
		initapp()
	}
	huma.Get(app.api, path, handler)
}

func RegisterAction(name string, callback func()) {
	if app.isInit == false {
		initapp()
	}
	app.notifyActions = append(app.notifyActions, notifyItem{
		Name:   name,
		Action: callback,
	})
}

func Init(port int) (err error) {
	app.port = port

	if app.window, err = walk.NewMainWindow(); err != nil {
		return errors.Wrapf(err, "Failed to create main window")
	}

	// if ins.Icon, err = walk.Resources.Icon("syncthing.ico"); err != nil {
	// 	return nil, errors.Wrapf(err, "Failed to load icon")
	// }
	app.notify, err = walk.NewNotifyIcon(app.window)
	if err != nil {
		return errors.Wrapf(err, "Failed to create notify icon")
	}
	defer func() {
		if err != nil {
			app.notify.Dispose()
		}
	}()
	// if err = ins.Notify.SetIcon(ins.Icon); err != nil {
	// 	return nil, errors.Wrapf(err, "Failed to set icon for notify icon")
	// }
	if err = app.notify.SetVisible(true); err != nil {
		return errors.Wrapf(err, "Failed to set visible for notify icon")
	}

	// Setup Menus
	for _, item := range app.notifyActions {
		action := walk.NewAction()
		if err := action.SetText(item.Name); err != nil {
			return errors.Wrapf(err, "Failed to set text for quit action")
		}
		action.Triggered().Attach(item.Action)
		if err := app.notify.ContextMenu().Actions().Add(action); err != nil {
			return errors.Wrapf(err, "Failed to add action for %s", item.Name)
		}
	}

	return nil
}

func Run() {
	go func() {
		addr := fmt.Sprintf("127.0.0.1:%d", app.port)
		log.Infof("Listening on %s", addr)
		if err := http.ListenAndServe(addr, app.router); err != nil && err != http.ErrServerClosed {
			log.Errorf("Server error] %+v\n", err)
			Destroy()
		}
	}()

	app.window.Run()
}

func Destroy() {
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	// if err := ins.Server.Shutdown(ctx); err != nil {
	// 	log.WithError(err).Error("Failed to shutdown server")
	// }

	walk.App().Exit(0)
	app.notify.Dispose()
}
