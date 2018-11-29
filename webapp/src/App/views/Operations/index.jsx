import React, { PureComponent } from 'react'
import { Responsive as ResponsiveReactGridLayout } from 'react-grid-layout'
import { HotKeys } from 'react-hotkeys'
import OpSelector from '../../OpSelector'
import { AutoSizer } from 'react-virtualized'
import Item from './Item'
import 'react-grid-layout/css/styles.css'
import opspecNodeApiClient from '../../../core/clients/opspecNodeApi'
import contentStore from '../../../core/contentStore'
import uuidV4 from 'uuid/v4'
import { toast } from 'react-toastify'

const dragHandleClassName = 'dragHandle'
const CONTENT_STORE_KEY = 'operations'

export default class OperationsView extends PureComponent {
  static defaultProps = {
    className: 'layout',
    cols: { lg: 12, md: 10, sm: 6, xs: 4, xxs: 2 },
    rowHeight: 100
  };

  constructor (props) {
    super(props)

    this.state = contentStore.get({ key: CONTENT_STORE_KEY }) ||
      {
        layouts: {},
        items: []
      }
  }

  isItemStartable = (inputs, args) => Object.keys(inputs || []).length === Object.keys(args).length;

  addItem = ({ op, opRef }) => {
    const isStartable = this.isItemStartable(op.inputs, {})
    const item = {
      opRef,
      op,
      args: {},
      i: uuidV4(),
      x: (this.state.items.length * 2) % (12),
      y: 10000000000000, // puts it at the bottom
      w: 2,
      h: 2,
      isStartable
    }

    this.setState(
      prevState => ({
        items: [...prevState.items, item]
      })
    )
  };

  componentDidUpdate () {
    contentStore.set({ key: CONTENT_STORE_KEY, value: this.state })
  }

  handleLayoutChange = (layout, layouts) => {
    this.setState({ layouts })
  };

  deleteItem = (itemId) => {
    this.setState(
      prevState => ({ items: prevState.items.filter(item => item.i !== itemId) })
    )
  };

  toggleFullScreenItem = (itemId) => {
    this.setState(
      prevState => {
        const itemIndex = prevState.items.findIndex(item => item.i === itemId)
        const items = [...prevState.items]
        const item = prevState.items[itemIndex]
        items[itemIndex] = Object.assign({}, item, { isFullScreen: !item.isFullScreen })
        return { items }
      }
    )
  };

  updateItemConfiguration = (itemId, configuration) => {
    this.setState(
      prevState => {
        const itemIndex = prevState.items.findIndex(item => item.i === itemId)
        const items = [...prevState.items]
        const item = Object.assign({}, prevState.items[itemIndex], configuration)
        item.isStartable = this.isItemStartable(item.op.inputs, item.args)
        items[itemIndex] = item
        return { items }
      }
    )
  };

  selectAllItems = () => {
    this.setState({ isAllItemsSelected: true })
    return false
  };

  unSelectAllItems = () => {
    this.setState({ isAllItemsSelected: false })
  };

  startItem = (itemId) => {
    const item = this.state.items.find(item => item.i === itemId)

    if (!item.isStartable) {
      toast.error(`Unable to start ${item.name || item.opRef}; configuration required`)
      return
    }

    const args = Object.entries(item.op.inputs || [])
      .reduce((args, [name, param]) => {
        if (param.array) args[name] = { array: item.args[name] }
        if (param.boolean) args[name] = { boolean: item.args[name] }
        if (param.dir) args[name] = { dir: item.args[name] }
        if (param.file) args[name] = { file: item.args[name] }
        if (param.number) args[name] = { number: item.args[name] }
        if (param.object) args[name] = { object: item.args[name] }
        if (param.socket) args[name] = { socket: item.args[name] }
        if (param.string) args[name] = { string: item.args[name] }
        return args
      }, {})

    opspecNodeApiClient.op_start({
      args,
      op: {
        ref: item.opRef
      }
    })
      .then(opId => {
        this.updateItemConfiguration(itemId, { opId })
      })
      .catch(error => {
        toast.error(error.message)
      })
  };

  startAllItems = () => {
    if (!this.state.isAllItemsSelected) return
    this.state.items.forEach(item => this.startItem(item.i))
  };

  killItem = (itemId) => {
    const item = this.state.items.find(item => item.i === itemId)
    opspecNodeApiClient.op_kill({
      opId: item.opId
    })
      .then(() => {
        this.updateItemConfiguration(item.i, { isKillable: false })
      })
      .catch(error => {
        toast.error(error.message)
      })
  };

  killAllItems = () => {
    if (!this.state.isAllItemsSelected) return
    this.state.items.forEach(item => this.killItem(item.i))
  };

  render () {
    const fullScreenItem = this.state.items.find(item => item.isFullScreen)
    if (fullScreenItem) {
      return (<Item
        opId={fullScreenItem.opId}
        opRef={fullScreenItem.opRef}
        op={fullScreenItem.op}
        name={fullScreenItem.name}
        isFullScreen
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
          'startAllItems': 'enter'
        }}
        handlers={{
          'selectAllItems': this.selectAllItems,
          'killAllItems': this.killAllItems,
          'startAllItems': this.startAllItems
        }}
        onClick={this.unSelectAllItems}
      >
        <OpSelector
          onSelect={this.addItem}
        />
        <AutoSizer>
          {({ width }) =>
            <ResponsiveReactGridLayout
              width={width}
              // avoids creation of stacking context per item which causes dropdown from one item to render behind other items
              useCSSTransforms={false}
              layouts={this.state.layouts}
              onLayoutChange={this.handleLayoutChange}
              draggableHandle={`.${dragHandleClassName}`}
              // ensures grid is "selectable" so hotkeys will work within it
              style={{ width }}
            >
              {this.state.items.map(item =>
                <div
                  data-grid={item}
                  key={item.i}>
                  <Item
                    opId={item.opId}
                    opRef={item.opRef}
                    op={item.op}
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
