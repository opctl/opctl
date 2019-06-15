import CallEnded from './callEnded'
import ContainerExited from './containerExited'
import ContainerStarted from './containerStarted'
import ContainerStdErrWrittenTo from './containerStdErrWrittenTo'
import ContainerStdOutWrittenTo from './containerStdOutWrittenTo'
import OpEnded from './opEnded'
import OpErred from './opErred'
import OpStarted from './opStarted'
import ParallelCallEndedEvent from './parallelCallEnded'
import SerialCallEndedEvent from './serialCallEnded'

export default interface Event {
    callEnded?: CallEnded | null | undefined
    containerExited?: ContainerExited | null | undefined
    containerStarted?: ContainerStarted | null | undefined
    containerStdErrWrittenTo?: ContainerStdErrWrittenTo | null | undefined
    containerStdOutWrittenTo?: ContainerStdOutWrittenTo | null | undefined
    opEnded?: OpEnded | null | undefined
    opErred?: OpErred | null | undefined
    opStarted?: OpStarted | null | undefined
    timestamp: Date
    parallelCallEnded?: ParallelCallEndedEvent | null | undefined
    serialCallEnded?: SerialCallEndedEvent | null | undefined
}