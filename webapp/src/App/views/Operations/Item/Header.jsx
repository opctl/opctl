import React, { PureComponent } from 'react'
import { MdMoreVert } from 'react-icons/md'
import { Dropdown, DropdownItem, DropdownMenu, DropdownToggle } from 'reactstrap'

const dragHandleClassName = 'dragHandle'

export default class Header extends PureComponent {
  state = { dropdownOpen: false };

  toggleDropdown = () => {
    this.setState(prevState => ({ dropdownOpen: !prevState.dropdownOpen }))
  };

  render () {
    return (
      <div
        className={dragHandleClassName}
        onDoubleClick={this.props.onToggleFullScreen}
        style={{
          position: 'absolute',
          width: '100%',
          height: '37px',
          left: '0',
          top: '0',
          borderBottom: 'solid thin #ececec',
          cursor: this.props.isFullScreen ? 'pointer' : '-webkit-grab',
          wordBreak: 'break-all',
          verticalAlign: 'middle',
          textAlign: 'center',
          lineHeight: '37px',
          whiteSpace: 'nowrap'
        }}
      >
        <div style={{ overflow: 'hidden', marginLeft: '4px', width: 'calc(100% - 36px)', height: '37px' }}>
          {this.props.name}
        </div>
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
                onClick={this.props.onConfigure}
              >
                Configure
              </DropdownItem>
              <DropdownItem
                onClick={this.props.onDelete}
              >
                Delete (del)
              </DropdownItem>
            </DropdownMenu>
          </Dropdown>
        </div>
      </div>
    )
  }
}
