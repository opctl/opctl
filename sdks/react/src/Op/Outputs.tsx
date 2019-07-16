import React from 'react'
import ModelParam from '@opctl/sdk/lib/model/param'
import Output from './Output/index'

interface Props {
  apiBaseUrl: string
  outputs: { [name: string]: ModelParam }
  opRef: string
  scope: { [name: string]: any }
}

export default (
  {
    apiBaseUrl,
    outputs = {},
    opRef,
    scope = {}
  }: Props
) => {
  if (!outputs || Object.entries(outputs).length === 0) return (null)

  return (
    <div>
      <h2>Outputs:</h2>
      {
        Object.entries(outputs).map(([name, param]) =>
          <Output
            apiBaseUrl={apiBaseUrl}
            key={name}
            name={name}
            param={param}
            opRef={opRef}
            value={scope[name] || {}}
          />
        )
      }
    </div>
  )
}
