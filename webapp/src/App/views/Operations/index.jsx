import React, {Component} from 'react';
import {Responsive as ResponsiveReactGridLayout} from "react-grid-layout";
import {HotKeys} from 'react-hotkeys';
import PkgSelector from '../../PkgSelector';
import {AutoSizer} from 'react-virtualized';
import Item from './Item';
import 'react-grid-layout/css/styles.css';
import opspecNodeApiClient from "../../../utils/clients/opspecNodeApi";
import {toast} from "react-toastify";

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

export default class Operations extends Component {
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

  isItemStartable = (inputs, args) => Object.keys(inputs || []).length === Object.keys(args).length;

  addItem = ({pkg, pkgRef}) => {
    const isStartable = this.isItemStartable(pkg.inputs, {});
    const item = {
      pkgRef,
      pkg,
      args: {},
      i: "n" + this.state.newCounter,
      x: (this.state.items.length * 2) % (12),
      y: 10000000000000, // puts it at the bottom
      w: 2,
      h: 2,
      isStartable,
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
        const item = Object.assign({}, prevState.items[itemIndex], configuration);
        item.isStartable = this.isItemStartable(item.pkg.inputs, item.args);
        items[itemIndex] = item;
        return {items};
      }
    );
  };

  selectAllItems = () => {
    this.setState({isAllItemsSelected: true});
  };

  unSelectAllItems = () => {
    this.setState({isAllItemsSelected: false});
  };

  startItem = (itemId) => {
    const item = this.state.items.find(item => item.i === itemId);

    if (!item.isStartable) {
      toast.error(`Unable to start ${item.name || item.pkgRef}; configuration required`);
      return;
    }

    const args = Object.entries(item.pkg.inputs || [])
      .reduce((args, [name, param]) => {
        if (param.array) args[name] = {array: item.args[name]};
        if (param.dir) args[name] = {dir: item.args[name]};
        if (param.file) args[name] = {file: item.args[name]};
        if (param.number) args[name] = {number: item.args[name]};
        if (param.object) args[name] = {object: item.args[name]};
        if (param.socket) args[name] = {socket: item.args[name]};
        if (param.string) args[name] = {string: item.args[name]};
        return args;
      }, {});

    opspecNodeApiClient.op_start({
      args,
      pkg: {
        ref: item.pkgRef,
      }
    })
      .then(opId => {
        this.updateItemConfiguration(itemId, {opId});
      })
      .catch(error => {
        toast.error(error.message);
      });
  };

  startAllItems = () => {
    if (!this.state.isAllItemsSelected) return;
    this.state.items.forEach(item => this.startItem(item.i));
  };

  killItem = (itemId) => {
    const item = this.state.items.find(item => item.i === itemId);
    opspecNodeApiClient.op_kill({
      opId: item.opId
    })
      .then(() => {
        this.updateItemConfiguration(item.i, {isKillable: false});
      })
      .catch(error => {
        toast.error(error.message);
      });
  };

  killAllItems = () => {
    if (!this.state.isAllItemsSelected) return;
    this.state.items.forEach(item => this.killItem(item.i));
  };

  render() {
    const fullScreenItem = this.state.items.find(item => item.isFullScreen);
    if (fullScreenItem) {
      return (<Item
        opId={fullScreenItem.opId}
        pkgRef={fullScreenItem.pkgRef}
        pkg={fullScreenItem.pkg}
        name={fullScreenItem.name}
        isFullScreen={true}
        isStartable={fullScreenItem.isStartable}
        onDelete={this.deleteItem.bind(this, fullScreenItem.i)}
        onStart={this.startItem.bind(this, fullScreenItem.i)}
        onKill={this.killItem.bind(this, fullScreenItem.i)}
        onToggleFullScreen={this.toggleFullScreenItem.bind(this, fullScreenItem.i)}
        args={fullScreenItem.args}
        onConfigured={this.updateItemConfiguration.bind(this, fullScreenItem.i)}
      />)
    }

    return (
      <HotKeys
        keyMap={{
          'selectAllItems': 'ctrl+a',
          'killAllItems': 'ctrl+c',
          'startAllItems': 'enter',
        }}
        handlers={{
          'selectAllItems': this.selectAllItems,
          'killAllItems': this.killAllItems,
          'startAllItems': this.startAllItems,
        }}
        style={{height: '100%'}}
        onClick={this.unSelectAllItems}
      >
        <PkgSelector
          onSelect={this.addItem}
        />
        <AutoSizer>
          {({width}) =>
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
                    name={item.name}
                    isAllItemsSelected={this.state.isAllItemsSelected}
                    isStartable={item.isStartable}
                    onDelete={this.deleteItem.bind(this, item.i)}
                    onKill={this.killItem.bind(this, item.i)}
                    onStart={this.startItem.bind(this, item.i)}
                    onToggleFullScreen={this.toggleFullScreenItem.bind(this, item.i)}
                    args={item.args}
                    onConfigured={this.updateItemConfiguration.bind(this, item.i)}
                  />
                </div>
              )}
            </ResponsiveReactGridLayout>
          }
        </AutoSizer>
      </HotKeys>
    )
  }
}
