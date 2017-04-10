package backend

import "github.com/laincloud/streamrouter/model"

type Backend interface {
	Check() error
	Reload()
	RenderStreamFiles([]model.StreamApp) error
}
