package connection

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"github.com/koer/koer-module/pkg/config"
	"google.golang.org/api/option"
)

func NewFirebaseApp(ctx context.Context, cfg config.FirebaseConfig) (*firebase.App, error) {
	var opts []option.ClientOption
	if cfg.CredentialsFile != "" {
		opts = append(opts, option.WithCredentialsFile(cfg.CredentialsFile))
	}
	app, err := firebase.NewApp(ctx, &firebase.Config{ProjectID: cfg.ProjectID}, opts...)
	if err != nil {
		return nil, fmt.Errorf("initializing firebase app: %w", err)
	}
	return app, nil
}
