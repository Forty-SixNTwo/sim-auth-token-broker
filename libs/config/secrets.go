package config

import (
	"context"
	"fmt"
	"os"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

func fetchSecrets() error {
	ctx := context.Background()
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("secretmanager.NewClient: %w", err)
	}
	defer client.Close()

	for _, env := range os.Environ() {
		parts := strings.SplitN(env, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, val := parts[0], parts[1]
		if strings.HasPrefix(val, "projects/") && strings.Contains(val, "/secrets/") {
			req := &secretmanagerpb.AccessSecretVersionRequest{Name: val}
			resp, err := client.AccessSecretVersion(ctx, req)
			if err != nil {
				return fmt.Errorf("failed to access secret version for %s: %w", key, err)
			}
			secretData := string(resp.Payload.Data)
			err = os.Setenv(key, secretData)
			if err != nil {
				return fmt.Errorf("failed to set env %s: %w", key, err)
			}
		}
	}
	return nil
}
