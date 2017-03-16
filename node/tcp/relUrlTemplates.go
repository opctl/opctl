package tcp

/* resources */
const (
	// resource: a single event stream
	eventsStreamRelUrlTemplate string = "/events/stream"

	// resource: a single liveness
	livenessRelUrlTemplate string = "/liveness"

	// resource: all instance kills
	opKillsRelUrlTemplate string = "/ops/kills"

	// resource: all instance starts
	opStartsRelUrlTemplate string = "/ops/starts"
)

/* use cases */
const (
	getLivenessRelUrlTemplate string = livenessRelUrlTemplate
	getPubSubRelUrlTemplate   string = eventsStreamRelUrlTemplate
	killOpRelUrlTemplate      string = opKillsRelUrlTemplate
	startOpRelUrlTemplate     string = opStartsRelUrlTemplate
)
