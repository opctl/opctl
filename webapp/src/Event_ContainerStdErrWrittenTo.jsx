import React from 'react';
import Ansi from 'ansi-to-react';

export default function EventContainerStdErrWrittenTo(props) {
  return (
    <div>
      <Ansi>
        {atob(props.containerStdErrWrittenTo.data)}
      </Ansi>
    </div>
  );
}
