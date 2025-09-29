package terraform

import (
	"autohost-cli/internal/ports"
	"context"
)

type Adapter struct{}

func New() *Adapter { return &Adapter{} }

func (a *Adapter) Install(ctx context.Context) error { return Install(ctx) }
func (a *Adapter) ApplySplitDNS(ctx context.Context, workProfile string, cfg ports.SplitDNSConfig) error {
	return ApplySplitDNS(ctx, workProfile, cfg)
}
