package retry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

func WithBackoff(ctx context.Context, operation func() error) error {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = 10 * time.Second // stop retrying after 10 seconds

	if err := backoff.Retry(operation, backoff.WithContext(expBackoff, ctx)); err != nil {
		return err
	}

	return nil
}
