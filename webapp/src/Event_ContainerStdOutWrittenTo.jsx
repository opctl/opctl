import React from 'react';
import Ansi from 'ansi-to-react';


export default function EventContainerStdOutWrittenTo(props) {
  return (
    <div>
      <Ansi>
        {atob(props.containerStdOutWrittenTo.data)}
      </Ansi>
    </div>
  );
}
