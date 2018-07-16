package bridge

import (
	"keep/docker"

	"github.com/docker/docker/client"
	"golang.org/x/net/context"
)

type Bridge struct {
	Docker *client.Client
	Ctx    *context.Context
}

func New(docker *client.Client, ctx *context.Context) *Bridge {
	return &Bridge{Docker: docker, Ctx: ctx}
}
func (b *Bridge) StartMonitor() {
	docker.MonitorStats(b.Docker, b.Ctx)
}
