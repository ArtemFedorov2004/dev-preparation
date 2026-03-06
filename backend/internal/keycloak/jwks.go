package keycloak

import (
	"context"
	"fmt"
	"time"

	"github.com/MicahParks/keyfunc/v3"
)

type Config struct {
	RealmURL        string
	RefreshInterval time.Duration
}

func NewJWKS(ctx context.Context, cfg Config) (keyfunc.Keyfunc, error) {
	jwksURL := cfg.RealmURL + "/protocol/openid-connect/certs"

	k, err := keyfunc.NewDefaultCtx(ctx, []string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("create jwks keyfunc: %w", err)
	}

	return k, nil
}
