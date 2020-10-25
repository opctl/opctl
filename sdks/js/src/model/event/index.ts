import CallEnded from './callEnded'
import CallKilled from './callKilled'
import ContainerExited from './containerExited'
import ContainerStarted from './containerStarted'
import ContainerStdErrWrittenTo from './containerStdErrWrittenTo'
import ContainerStdOutWrittenTo from './containerStdOutWrittenTo'
import OpEnded from './opEnded'
import OpErred from './opErred'
import OpStarted from './opStarted'
import ParallelCallEnded from './parallelCallEnded'
import SerialCallEnded from './serialCallEnded'

export default interface Event {
    callEnded?: CallEnded | null | undefined
    callKilled?: CallKilled | null | undefined
    containerExited?: ContainerExited | null | undefined
    containerStarted?: ContainerStarted | null | undefined
    containerStdErrWrittenTo?: ContainerStdErrWrittenTo | null | undefined
    containerStdOutWrittenTo?: ContainerStdOutWrittenTo | null | undefined
    opEnded?: OpEnded | null | undefined
    opErred?: OpErred | null | undefined
    opStarted?: OpStarted | null | undefined
    timestamp: Date
    parallelCallEnded?: ParallelCallEnded | null | undefined
    serialCallEnded?: SerialCallEnded | null | undefined
}