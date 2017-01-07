package tcp

/* resources */
const (
	// resource: a single event-stream
	eventBusRelUrlTemplate string = "/event-stream"

	// resource: a single liveness
	livenessRelUrlTemplate string = "/liveness"

	// resource: all instance kills
	opKillsRelUrlTemplate string = "/instances/kills"

	// resource: all instance starts
	opStartsRelUrlTemplate string = "/instances/starts"
)

/* use cases */
const (
	getLivenessRelUrlTemplate string = livenessRelUrlTemplate
	getEventBusRelUrlTemplate string = eventBusRelUrlTemplate
	killOpRelUrlTemplate      string = opKillsRelUrlTemplate
	startOpRelUrlTemplate     string = opStartsRelUrlTemplate
)
