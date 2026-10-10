package observability

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/truongle2004/mercato/pkg/config"
)

func TestSetupTracing_Disabled(t *testing.T) {
	shutdown, err := SetupTracing(context.Background(), &config.Schema{})

	require.NoError(t, err)
	require.NotNil(t, shutdown)
	require.NoError(t, shutdown(context.Background()))
}
