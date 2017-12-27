import React from 'react';
import Output from './Output/index';

export default function Outputs({params = {}, pkgRef, values = {}}) {
  return (
    <div>
      <h2>Outputs:</h2>
      {
        Object.entries(params).map(([name, param]) =>
          <Output
            key={name}
            name={name}
            param={param}
            pkgRef={pkgRef}
            value={values[name] || {}}
          />
        )
      }
    </div>
  );
}
