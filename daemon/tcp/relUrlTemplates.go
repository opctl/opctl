package tcp

/* resources */
const (
	// resource: a single event-stream
	pubSubRelUrlTemplate string = "/event-stream"

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
	getPubSubRelUrlTemplate   string = pubSubRelUrlTemplate
	killOpRelUrlTemplate      string = opKillsRelUrlTemplate
	startOpRelUrlTemplate     string = opStartsRelUrlTemplate
)
