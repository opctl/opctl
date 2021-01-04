import React, { useState } from 'react'
import AutoComplete from '../AutoComplete'

type CallType = {
  name: string
  id: string
}

const callTypes: CallType[] = [
  {
    name: 'Op',
    id: '07e87a89-2b2f-47d3-91fd-d8c1d1178e92'
  },
  {
    name: 'Container',
    id: '13dae8ef-e09d-4a6c-843e-de8e64f817d8'
  },
  {
    name: 'Serial',
    id: 'efc3724e-6f67-4518-8fbf-2141f7a81d72'
  },
  {
    name: 'Serial Loop',
    id: '99bc32c6-6cb2-4d52-a677-b15dd6ab57ad'
  },
  {
    name: 'Parallel',
    id: '7856bc12-540b-47ff-bebd-239a141fb42c'
  },
  {
    name: 'Parallel Loop',
    id: 'b2f43f0f-5fc0-47f4-a972-effb7e291237'
  }
]

interface Props {
  onChange: (value: string) => any
}

export default function CallTypeInput(
  {
    onChange
  }: Props
) {
  const [suggestions, setSuggestions] = useState([] as any[])

  return (
    <AutoComplete<CallType>
      onSearch={value => {
        const inputValue = value.trim().toLowerCase()
        const inputLength = inputValue.length

        setSuggestions(inputLength === 0 ? [] : callTypes.filter(callType =>
          callType.name.toLowerCase().slice(0, inputLength) === inputValue
        ))
      }}
      getValue={callType => callType.name}
      onSelect={callType => onChange(callType.name)}
      options={suggestions}
      placeholder='Search calls...'
      render={suggestion => (
        <div>
          {suggestion.name}
        </div>
      )}
    />
  )
}