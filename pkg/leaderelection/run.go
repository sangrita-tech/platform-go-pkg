package leaderelection

import (
	"context"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func (e *Elector) Run(ctx context.Context) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		lock := &resourcelock.LeaseLock{
			LeaseMeta: metav1.ObjectMeta{
				Name:      e.cfg.LeaseName,
				Namespace: e.cfg.LeaseNamespace,
			},
			Client: e.clientset.CoordinationV1(),
			LockConfig: resourcelock.ResourceLockConfig{
				Identity: e.identity,
			},
		}

		electionCtx, cancelElection := context.WithCancel(ctx)

		lec := leaderelection.LeaderElectionConfig{
			Lock:            lock,
			LeaseDuration:   e.cfg.LeaseDuration,
			RenewDeadline:   e.cfg.RenewDeadline,
			RetryPeriod:     e.cfg.RetryPeriod,
			ReleaseOnCancel: true,
			Name:            e.cfg.LeaseName,
			Callbacks: leaderelection.LeaderCallbacks{
				OnStartedLeading: func(leaderCtx context.Context) {
					if e.cb.OnNewLeader != nil {
						e.cb.OnNewLeader(e.identity)
					}
					if e.cb.OnStartedLeading != nil {
						e.cb.OnStartedLeading(leaderCtx)
					}
				},
				OnStoppedLeading: func() {
					cancelElection()
					if e.cb.OnStoppedLeading != nil {
						e.cb.OnStoppedLeading()
					}
				},
				OnNewLeader: func(id string) {
					if e.cb.OnNewLeader != nil && id != "" {
						e.cb.OnNewLeader(id)
					}
				},
			},
		}

		le, err := leaderelection.NewLeaderElector(lec)
		if err != nil {
			cancelElection()
			return err
		}

		le.Run(electionCtx)
		cancelElection()

		if ctx.Err() != nil {
			return nil
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(e.cfg.RetryPeriod):
		}
	}
}
