import React, {Component} from 'react';

/**
 * A higher order component for components whose content never changes thus optimized by always returning false from
 * shouldComponentUpdate
 * @param WrappedComponent
 * @returns {{new(): NeverReRender}}
 */
export default (WrappedComponent) => {
  return class NeverReRender extends Component {
    shouldComponentUpdate() {
      return false;
    }

    render() {
      return <WrappedComponent {...this.props} />
    }
  }
}
