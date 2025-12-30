package leaderelection_test

import (
	"context"
	"testing"
	"time"

	"github.com/sangrita-tech/platform-go-pkg/pkg/leaderelection"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func Test_New_ValidatesClientsetAndContextCanceled(t *testing.T) {
	t.Parallel()

	cfg := &leaderelection.Config{
		LeaseDuration: 60 * time.Second,
		RenewDeadline: 20 * time.Second,
		RetryPeriod:   5 * time.Second,
	}

	e, err := leaderelection.New(cfg, leaderelection.Callbacks{}, nil)
	require.Error(t, err)
	require.Nil(t, e)

	testEnv := &envtest.Environment{}
	restCfg, err := testEnv.Start()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = testEnv.Stop()
	})

	clientset, err := kubernetes.NewForConfig(restCfg)
	require.NoError(t, err)

	cfg2 := &leaderelection.Config{
		LeaseName:      "lease",
		LeaseNamespace: "default",
		Identity:       "x",
		LeaseDuration:  4 * time.Second,
		RenewDeadline:  3 * time.Second,
		RetryPeriod:    1 * time.Second,
	}
	e2, err := leaderelection.New(cfg2, leaderelection.Callbacks{}, clientset)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = e2.Run(ctx)
	require.ErrorIs(t, err, context.Canceled)
}
