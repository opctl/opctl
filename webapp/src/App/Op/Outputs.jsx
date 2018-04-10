import React from 'react'
import Output from './Output/index'

export default ({params = {}, opRef, values = {}}) => {
  if (!params || Object.entries(params).length === 0) return (null)

  return (
    <div>
      <h2>Outputs:</h2>
      {
        Object.entries(params).map(([name, param]) =>
          <Output
            key={name}
            name={name}
            param={param}
            opRef={opRef}
            value={values[name] || {}}
          />
        )
      }
    </div>
  )
}
