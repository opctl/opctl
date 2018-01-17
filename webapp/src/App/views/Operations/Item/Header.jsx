import React, {Component} from 'react';
import MdMoreVert from 'react-icons/lib/md/more-vert';
import {Dropdown, DropdownItem, DropdownMenu, DropdownToggle} from 'reactstrap';

const dragHandleClassName = 'dragHandle';

export default class Header extends Component {
  state = {dropdownOpen: false};

  toggleDropdown = () => {
    this.setState(prevState => ({dropdownOpen: !prevState.dropdownOpen}));
  };

  render() {
    return (
      <div
        className={dragHandleClassName}
        style={{
          position: 'absolute',
          width: '100%',
          height: '37px',
          top: '0',
          borderBottom: 'solid thin #ececec',
          cursor: '-webkit-grab',
          wordBreak: 'break-all',
          verticalAlign: 'middle',
          textAlign: 'center',
          lineHeight: '37px',
          whiteSpace: 'nowrap'
        }}
      >
        <div style={{overflow: 'hidden', width: 'calc(100% - 36px)', height: '37px'}}>
          {this.props.pkgRef}
        </div>
        <div style={{right: 0, top: 0, position: 'absolute'}}>
          <Dropdown
            isOpen={this.state.dropdownOpen}
            toggle={this.toggleDropdown}
          >
            <DropdownToggle tag="div">
              <MdMoreVert
                style={{
                  transform: 'rotate(-90deg) translateY(-50%)',
                  cursor: 'pointer'
                }}
              />
            </DropdownToggle>
            <DropdownMenu right>
              {
                this.props.isKillable
                  ?
                  <DropdownItem onClick={this.props.onKill}>
                    Kill
                  </DropdownItem>
                  :
                  <DropdownItem onClick={this.props.onStart} disabled={!this.props.isStartable}>
                    Start
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
                Delete
              </DropdownItem>
            </DropdownMenu>
          </Dropdown>
        </div>
      </div>
    );
  }
}
