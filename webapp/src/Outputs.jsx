import React from 'react';
import Output from './Output';

export default function Outputs({value}) {

  const outputs = Object.entries(value || {}).map(([name, output]) => (
    <Output
      name={name}
      output={output}
      key={name}
    />
  ));

  return (
    <div>
      <h2>Outputs:</h2>
      {outputs}
    </div>
  );
}
