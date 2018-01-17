import React, { Component } from 'react';
import { Responsive as ResponsiveReactGridLayout } from "react-grid-layout";
import PkgSelector from '../../PkgSelector';
import {AutoSizer} from 'react-virtualized';
import Item from './Item';
import 'react-grid-layout/css/styles.css';

const dragHandleClassName = 'dragHandle';

function getStateFromLS(key) {
  if (global.localStorage) {
    return JSON.parse(global.localStorage.getItem("state")) || null;
  }
  return null;
}

function saveStateToLS(state) {
  if (global.localStorage) {
    global.localStorage.setItem(
      "state",
      JSON.stringify(state)
    );
  }
}

export default class Spaces extends Component {
  static defaultProps = {
    className: "layout",
    cols: {lg: 12, md: 10, sm: 6, xs: 4, xxs: 2},
    rowHeight: 100
  };

  constructor(props) {
    super(props);

    this.state =
      getStateFromLS()
      ||
      {
        layouts: {},
        items: [],
        newCounter: 0
      };
  }

  addItem = ({pkg, pkgRef}) => {
    const item = {
      pkgRef,
      pkg,
      args: {},
      i: "n" + this.state.newCounter,
      x: (this.state.items.length * 2) % (12),
      y: 10000000000000, // puts it at the bottom
      w: 2,
      h: 2
    };

    this.setState(
      prevState => ({
        items: [...prevState.items, item],
        newCounter: prevState.newCounter + 1
      })
    );
  };

  componentDidUpdate() {
    saveStateToLS(this.state);
  }

  handleLayoutChange = (layout, layouts) => {
    this.setState({layouts});
  };

  deleteItem = (itemId) => {
    this.setState(
      prevState => ({items: prevState.items.filter(item => item.i !== itemId)})
    );
  };

  toggleFullScreenItem = (itemId) => {
    this.setState(
      prevState => {
        const itemIndex = prevState.items.findIndex(item => item.i === itemId);
        const items = [...prevState.items];
        const item = prevState.items[itemIndex];
        items[itemIndex] = Object.assign({}, item, {isFullScreen: !item.isFullScreen});
        return {items};
      }
    );
  };

  updateItemConfiguration = (itemId, configuration) => {
    this.setState(
      prevState => {
        const itemIndex = prevState.items.findIndex(item => item.i === itemId);
        const items = [...prevState.items];
        const item = prevState.items[itemIndex];
        items[itemIndex] = Object.assign({}, item, configuration);
        return {items};
      }
    );
  };

  render() {
    const fullScreenItem = this.state.items.find(item => item.isFullScreen);
    return (
      <div style={{height: '100%'}}>
        {
          fullScreenItem
            ?
            <Item
              opId={fullScreenItem.opId}
              pkgRef={fullScreenItem.pkgRef}
              pkg={fullScreenItem.pkg}
              isFullScreen={true}
              onDelete={this.deleteItem.bind(this, fullScreenItem.i)}
              onToggleFullScreen={this.toggleFullScreenItem.bind(this, fullScreenItem.i)}
              args={fullScreenItem.args}
              onConfigured={this.updateItemConfiguration.bind(this, fullScreenItem.i)}
            />
            :
            <div>
              <PkgSelector
                onSelect={this.addItem}
              />
              <AutoSizer>
                {({height, width}) =>
                  <ResponsiveReactGridLayout
                    width={width}
                    // avoids creation of stacking context per item which causes dropdown from one item to render behind other items
                    useCSSTransforms={false}
                    layouts={this.state.layouts}
                    onLayoutChange={this.handleLayoutChange}
                    draggableHandle={`.${dragHandleClassName}`}
                  >
                    {this.state.items.map(item =>
                      <div
                        data-grid={item}
                        key={item.i}>
                        <Item
                          opId={item.opId}
                          pkgRef={item.pkgRef}
                          pkg={item.pkg}
                          onDelete={this.deleteItem.bind(this, item.i)}
                          onToggleFullScreen={this.toggleFullScreenItem.bind(this, item.i)}
                          args={item.args}
                          onConfigured={this.updateItemConfiguration.bind(this, item.i)}
                        />
                      </div>
                    )}
                  </ResponsiveReactGridLayout>
                }
              </AutoSizer>
            </div>
        }
      </div>
    )
  }
}
