package leaderelection

import "context"

type Callbacks struct {
	OnStartedLeading func(leaderCtx context.Context)
	OnStoppedLeading func()
	OnNewLeader      func(identity string)
}
