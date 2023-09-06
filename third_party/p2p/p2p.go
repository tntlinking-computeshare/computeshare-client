package p2p

import "github.com/google/wire"

var ProviderSet = wire.NewSet(RunDaemon)
