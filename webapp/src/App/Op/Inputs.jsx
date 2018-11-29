import React from 'react'
import Input from './Input/index'

export default ({ value, onInvalid, onValid, opRef, values }) => {
  if (!value || Object.entries(value).length === 0) return (null)

  return (
    <div>
      <h2>Inputs</h2>
      {Object.entries(value).map(([name, input]) =>
        <Input
          onInvalid={() => (onInvalid(name))}
          onValid={value => (onValid(name, value))}
          name={name}
          opRef={opRef}
          input={input}
          key={name}
          value={values[name] || null}
        />
      )}
    </div>
  )
}
