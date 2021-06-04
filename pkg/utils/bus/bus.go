package bus

import (
	evbus "github.com/asaskevich/EventBus"
)

var (
	bus                = evbus.New()
	Subscribe          = bus.Subscribe
	SubscribeOnce      = bus.SubscribeOnce
	HasCallback        = bus.HasCallback
	Unsubscribe        = bus.Unsubscribe
	Publish            = bus.Publish
	SubscribeAsync     = bus.SubscribeAsync
	SubscribeOnceAsync = bus.SubscribeOnceAsync
	WaitAsync          = bus.WaitAsync
)
