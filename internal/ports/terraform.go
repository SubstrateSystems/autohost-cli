package ports

import "context"

type Terraform interface {
	Install(ctx context.Context) error
	ApplySplitDNS(ctx context.Context, workProfile string, cfg SplitDNSConfig) error
}

type SplitDNSConfig struct {
	MagicDNS         bool
	SearchPaths      []string
	SplitNameservers map[string][]string
}
