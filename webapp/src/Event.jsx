import React from 'react';
import EventContainerExited from './Event_ContainerExited';
import EventContainerStarted from './Event_ContainerStarted';
import EventContainerStdErrWrittenTo from './Event_ContainerStdErrWrittenTo';
import EventContainerStdOutWrittenTo from './Event_ContainerStdOutWrittenTo';
import EventOpEnded from './Event_OpEnded';
import EventOpErred from './Event_OpErred';
import EventOpStarted from './Event_OpStarted';

export default function Event(props) {
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
  }
  return null
}
