import React from 'react';
import Ansi from 'ansi-to-react';


export default ({
                  containerStdOutWrittenTo
                }) => {
  return (
    <div>
      <Ansi>
        {atob(containerStdOutWrittenTo.data)}
      </Ansi>
    </div>
  );
}
