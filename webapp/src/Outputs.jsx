import React from 'react';
import Output from './Output';

export default function Outputs({params = {}, values = {}}) {
  return (
    <div>
      <h2>Outputs:</h2>
      {
        Object.entries(params).map(([name, param]) =>
          <Output
            key={name}
            name={name}
            param={param}
            value={values[name] || {}}
          />
        )
      }
    </div>
  );
}
