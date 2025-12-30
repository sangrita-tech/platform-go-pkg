package leaderelection_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/sangrita-tech/platform-go-pkg/pkg/leaderelection"
	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

func Test_Run_TwoElectors_LeadershipTransfers(t *testing.T) {
	t.Parallel()

	testEnv := &envtest.Environment{}
	restCfg, err := testEnv.Start()
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = testEnv.Stop()
	})

	clientset, err := kubernetes.NewForConfig(restCfg)
	require.NoError(t, err)

	leaseName := "e2e-lease"
	leaseNS := "default"

	leaderCh := make(chan string, 20)

	mkCallbacks := func(id string) leaderelection.Callbacks {
		return leaderelection.Callbacks{
			OnStartedLeading: func(ctx context.Context) {
				select {
				case leaderCh <- id:
				default:
				}
			},
			OnNewLeader: func(identity string) {
				if identity != "" {
					select {
					case leaderCh <- identity:
					default:
					}
				}
			},
		}
	}

	cfg1 := &leaderelection.Config{
		LeaseName:      leaseName,
		LeaseNamespace: leaseNS,
		Identity:       "id-1",
		LeaseDuration:  4 * time.Second,
		RenewDeadline:  3 * time.Second,
		RetryPeriod:    1 * time.Second,
	}
	cfg2 := &leaderelection.Config{
		LeaseName:      leaseName,
		LeaseNamespace: leaseNS,
		Identity:       "id-2",
		LeaseDuration:  4 * time.Second,
		RenewDeadline:  3 * time.Second,
		RetryPeriod:    1 * time.Second,
	}

	e1, err := leaderelection.New(cfg1, mkCallbacks(cfg1.Identity), clientset)
	require.NoError(t, err)
	e2, err := leaderelection.New(cfg2, mkCallbacks(cfg2.Identity), clientset)
	require.NoError(t, err)

	rootCtx, rootCancel := context.WithCancel(context.Background())
	t.Cleanup(rootCancel)

	ctx1, cancel1 := context.WithCancel(rootCtx)
	ctx2, cancel2 := context.WithCancel(rootCtx)
	defer cancel1()
	defer cancel2()

	errCh := make(chan error, 2)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		errCh <- e1.Run(ctx1)
	}()
	go func() {
		defer wg.Done()
		errCh <- e2.Run(ctx2)
	}()

	var firstLeader string
	require.Eventually(t, func() bool {
		select {
		case v := <-leaderCh:
			if v == "id-1" || v == "id-2" {
				firstLeader = v
				return true
			}
			return false
		default:
			return false
		}
	}, 15*time.Second, 50*time.Millisecond)

	if firstLeader == "id-1" {
		cancel1()
	} else {
		cancel2()
	}

	var secondLeader string
	require.Eventually(t, func() bool {
		select {
		case v := <-leaderCh:
			if v == "id-1" || v == "id-2" {
				if v != firstLeader {
					secondLeader = v
					return true
				}
			}
			return false
		default:
			return false
		}
	}, 20*time.Second, 50*time.Millisecond)

	require.NotEqual(t, firstLeader, secondLeader)

	rootCancel()
	wg.Wait()

	ea := <-errCh
	eb := <-errCh
	require.NoError(t, ea)
	require.NoError(t, eb)
}
