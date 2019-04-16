import React, { PureComponent } from 'react'
import { MdMoreVert } from 'react-icons/md'
import { Dropdown, DropdownItem, DropdownMenu, DropdownToggle } from 'reactstrap'

interface Props {
  isFullScreen?: boolean | null | undefined
  isKillable: boolean
  isStartable: boolean
  name: string
  onDelete?: () => any | undefined | null
  onKill: () => any
  onStart: () => any
  onToggleFullScreen: () => any
}

interface State {
  dropdownOpen: boolean
}

export default class Header extends PureComponent<Props, State> {
  state = { dropdownOpen: false };

  toggleDropdown = () => {
    this.setState(prevState => ({ dropdownOpen: !prevState.dropdownOpen }))
  };

  render() {
    return (
      <div
        style={{
          width: '100%',
          height: '37px',
          wordBreak: 'break-all',
          verticalAlign: 'middle',
          lineHeight: '37px',
          whiteSpace: 'nowrap',
          position: 'relative'
        }}
      >
        <div style={{
          overflow: 'hidden',
          marginLeft: '4px',
          width: 'calc(100% - 36px)',
          height: '37px'
        }}>
          {this.props.name}
        </div>
        <div>
          <div style={{ right: 0, top: 0, position: 'absolute', cursor: 'pointer' }}>
            <Dropdown
              isOpen={this.state.dropdownOpen}
              toggle={this.toggleDropdown}
            >
              <DropdownToggle tag='div'>
                <MdMoreVert
                  style={{
                    transform: 'rotate(-90deg) translateY(-50%)'
                  }}
                />
              </DropdownToggle>
              <DropdownMenu right>
                {
                  this.props.isKillable
                    ? <DropdownItem onClick={this.props.onKill}>
                      Kill (ctrl+c)
                  </DropdownItem>
                    : <DropdownItem onClick={this.props.onStart} disabled={!this.props.isStartable}>
                      Start (enter)
                  </DropdownItem>
                }
                <DropdownItem
                  onClick={this.props.onToggleFullScreen}
                >
                  Toggle Full Screen
                </DropdownItem>
                {
                  this.props.onDelete
                    ? <DropdownItem
                      onClick={this.props.onDelete}
                    >
                      Delete (del)
                  </DropdownItem>
                    : null
                }
              </DropdownMenu>
            </Dropdown>
          </div>
        </div>
      </div>
    )
  }
}
