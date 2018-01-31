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
import withNoReRenders from '../../hocs/withNoReRenders';

const Event = ({event, style}) => {
  // delegate to typed event component
  let component;
  if (event.containerExited) {
    component = <EventContainerExited containerExited={event.containerExited} timestamp={event.timestamp}/>;
  } else if (event.containerStarted) {
    component = <EventContainerStarted containerStarted={event.containerStarted} timestamp={event.timestamp}/>;
  } else if (event.containerStdErrWrittenTo) {
    component = <EventContainerStdErrWrittenTo containerStdErrWrittenTo={event.containerStdErrWrittenTo}/>;
  } else if (event.containerStdOutWrittenTo) {
    component = <EventContainerStdOutWrittenTo containerStdOutWrittenTo={event.containerStdOutWrittenTo}/>;
  } else if (event.opEnded) {
    component = <EventOpEnded opEnded={event.opEnded} timestamp={event.timestamp}/>;
  } else if (event.opErred) {
    component = <EventOpErred opErred={event.opErred} timestamp={event.timestamp}/>;
  } else if (event.opStarted) {
    component = <EventOpStarted opStarted={event.opStarted} timestamp={event.timestamp}/>;
  } else if (event.parallelCallEnded) {
    component = <EventParallelCallEnded parallelCallEnded={event.parallelCallEnded} timestamp={event.timestamp}/>;
  } else if (event.serialCallEnded) {
    component = <EventSerialCallEnded serialCallEnded={event.serialCallEnded} timestamp={event.timestamp}/>;
  } else {
    component = null;
  }
  return <div style={style}>
    {component}
  </div>
};

export default withNoReRenders(Event);
