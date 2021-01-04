import React, { useState } from 'react'
import brandColors from '../../brandColors'
import Logo from '../Logo'
import { ReactComponent as MenuIcon } from '../../icons/Menu.svg'
import WorkspaceExplorer from '../WorkspaceExplorer'
import WindowPane from '../WindowPane'
import { WindowProvider } from '../WindowContext'

export default function Workspace() {
  const [isSidePanelOpen, setIsSidePanelOpen] = useState(false)

  return (
    <div
      style={{
        display: 'flex',
        width: '100%',
        height: '100%',
        flexDirection: 'column'
      }}
    >
      <nav
        style={{
          flexShrink: 0,
          display: 'flex',
          alignItems: 'center',
          width: '100%',
          paddingLeft: '1.5rem',
          paddingRight: '1.5rem',
          height: '3rem',
          borderBottom: `solid .2rem ${brandColors.reallyLightGray}`
        }}
      >
        <MenuIcon
          onClick={() => setIsSidePanelOpen(isSidePanelOpen => !isSidePanelOpen)}
          style={{
            cursor: 'pointer',
            fill: brandColors.active
          }}
        />
        <Logo
          style={{
            marginLeft: '1.5rem'
          }}
        />
      </nav>
      <main
        style={{
          display: 'flex',
          width: '100%',
          height: '100%',
          overflow: 'hidden'
        }}
      >
        <WindowProvider>
          <aside
            style={{
              ...isSidePanelOpen
                ? { display: 'none' }
                : {},
              width: '20rem',
              borderRight: `solid .2rem ${brandColors.reallyLightGray}`
            }}
          >
            <WorkspaceExplorer />
          </aside>
          <WindowPane />
        </WindowProvider>
      </main>
    </div>
  )
}