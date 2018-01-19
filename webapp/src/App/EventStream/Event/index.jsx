import React from 'react';
import EventContainerExited from './ContainerExited';
import EventContainerStarted from './ContainerStarted';
import EventContainerStdErrWrittenTo from './ContainerStdErrWrittenTo';
import EventContainerStdOutWrittenTo from './ContainerStdOutWrittenTo';
import EventOpEnded from './OpEnded';
import EventOpErred from './OpErred';
import EventOpStarted from './OpStarted';
import EventSerialCallEnded from "./SerialCallEnded";
import EventParallelCallEnded from "./ParallelCallEnded";

export default (props) => {
  // delegate to component for event
  if (props.event.containerExited) {
    return (<EventContainerExited containerExited={props.event.containerExited} timestamp={props.event.timestamp}/>);
  } else if (props.event.containerStarted) {
    return (<EventContainerStarted containerStarted={props.event.containerStarted} timestamp={props.event.timestamp}/>);
  } else if (props.event.containerStdErrWrittenTo) {
    return (<EventContainerStdErrWrittenTo containerStdErrWrittenTo={props.event.containerStdErrWrittenTo}/>);
  } else if (props.event.containerStdOutWrittenTo) {
    return (<EventContainerStdOutWrittenTo containerStdOutWrittenTo={props.event.containerStdOutWrittenTo}/>);
  } else if (props.event.opEnded) {
    return (<EventOpEnded opEnded={props.event.opEnded} timestamp={props.event.timestamp}/>);
  } else if (props.event.opErred) {
    return (<EventOpErred opErred={props.event.opErred} timestamp={props.event.timestamp}/>);
  } else if (props.event.opStarted) {
    return (<EventOpStarted opStarted={props.event.opStarted} timestamp={props.event.timestamp}/>);
  } else if (props.event.parallelCallEnded) {
    return (<EventParallelCallEnded parallelCallEnded={props.event.parallelCallEnded} timestamp={props.event.timestamp}/>);
  } else if (props.event.serialCallEnded) {
    return (<EventSerialCallEnded serialCallEnded={props.event.serialCallEnded} timestamp={props.event.timestamp}/>);
  }
  return null
}
