import dataGet from './data/get'
import eventsStream from './events/stream'
import livenessGet from './liveness/get'
import opsStartsPost from './ops/starts/post'
import opsKillsPost from './ops/kills/post'

export { dataGet }
export { eventsStream as eventStreamGet }
export { livenessGet }
export { opsKillsPost as opKill }
export { opsStartsPost as opStart }