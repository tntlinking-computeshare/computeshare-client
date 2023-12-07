package biz

import (
	"github.com/google/wire"
	"github.com/mohaijiang/computeshare-client/internal/biz/vm"
)

// ProviderSet is biz providers.
var ProviderSet = wire.NewSet(
	NewP2pClient,
	vm.NewVirtManager,
	NewStorageProvider,
)
