import React from 'react'
import ReactDOM from 'react-dom'
import Op from './index'

it('renders without crashing', () => {
  /* arrange */
  const div = document.createElement('div')

  /* act/assert */
  ReactDOM.render(
    <Op
      args={{}}
      eventStore={{getStream: () => {}} as any}
      isKillable={false}
      isStartable={false}
      onConfigure={() => null}
      onDelete={() => null}
      onKill={() => null}
      onStart={() => null}
      op={{ description: '' }}
      opId=''
      opRef=''
    />,
    div
  )
})
