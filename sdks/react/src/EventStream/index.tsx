import React, { Component } from 'react'
import Event from './Event'
import ModelEvent from '@opctl/sdk/src/model/event'
import EventStore from '../eventStore'
import { EventFilter } from '@opctl/sdk/src/api/client/events/stream'
import {
  AutoSizer,
  CellMeasurer,
  CellMeasurerCache,
  List
} from 'react-virtualized'
import {css} from '@emotion/css'
import { MeasuredCellParent } from 'react-virtualized/dist/es/CellMeasurer'

interface Props {
  eventStore: EventStore
  filter: EventFilter
}

interface State {
  events: ModelEvent[]
  throttleInterval?: number | null | undefined
}

const throttleDuration = 200

class EventStream extends Component<Props, State> {
  state = {
    events: []
  }

  cache = new CellMeasurerCache({
    fixedWidth: true
  })
  eventStreamCloser?: () => void | null
  interval: NodeJS.Timeout | null = null
  isScrolledToEnd: boolean = true
  isThrottled: boolean = false
  lastComponentUpdate = new Date().getTime()
  list?: List | null
  mostRecentWidth = 0

  componentDidMount() {
    this.eventStreamCloser = this.props.eventStore.getStream(
      this.props.filter,
      (event: ModelEvent) => {
        this.setState(prevState => ({
          events: [...prevState.events, event]
        }))
      })

    // maintain an update interval so throttled renders get re-processed every throttleDuration
    this.interval = setInterval(() => {
      this.setState(() => ({
        throttleInterval: new Date().getTime()
      }))
    }, throttleDuration)
  }

  componentDidUpdate() {
    this.lastComponentUpdate = new Date().getTime()
  }

  componentWillUnmount() {
    this.eventStreamCloser && this.eventStreamCloser()

    this.interval && clearInterval(this.interval)
  }

  rowRenderer = (
    {
      index,
      key,
      parent,
      style
    }: {
      index: number,
      key: string,
      parent: MeasuredCellParent,
      style: React.CSSProperties
    }
  ) => {
    return (
      <CellMeasurer
        cache={this.cache}
        columnIndex={0}
        key={key}
        parent={parent}
        rowIndex={index}
        width={this.mostRecentWidth}
      >
        <div key={key} style={style}>
          <Event event={this.state.events[index]} />
        </div>
      </CellMeasurer>
    )
  };

  handleResize = (
    {
      width
    }: {
      width: number
    }
  ) => {
    if (this.mostRecentWidth !== width) {
      this.mostRecentWidth = width
      this.cache.clearAll()
      this.list && this.list.recomputeRowHeights()
    }
  };

  handleRowsRendered = (
    {
      stopIndex
    }: {
      stopIndex: number
    }
  ) => {
    this.isScrolledToEnd = this.state.events.length - stopIndex <= 1
  };

  render() {
    return (
      this.state.events.length !== 0
        ? <AutoSizer onResize={this.handleResize}>
          {({ height, width }) => {
            const optionalProps = {} as any
            if (this.isScrolledToEnd) {
              optionalProps.scrollToIndex = this.state.events.length
            }
            return (
              <List
                className={css({
                  backgroundColor: '#222222',
                  wordBreak: 'break-all',
                  overflowWrap: 'break-word',
                  height: '100%',
                  overflow: 'auto',
                  code: {
                    color: '#f1f1f1',
                    backgroundColor: '#222222'
                  }
                })}
                height={height}
                width={width}
                ref={ref => {
                  this.list = ref
                }}
                onRowsRendered={this.handleRowsRendered}
                deferredMeasurementCache={this.cache}
                rowHeight={this.cache.rowHeight}
                rowRenderer={this.rowRenderer}
                rowCount={this.state.events.length}
                {...optionalProps}
              />
            )
          }}
        </AutoSizer>
        : null
    )
  }

  shouldComponentUpdate(nextProps: any, nextState: any) {
    // there are typically so many events received it floods react & results in a completely unresponsive UI
    // we throttle renders to avoid this
    const wasThrottled = this.isThrottled
    const isNewEvents = this.state.events.length < nextState.events.length
    this.isThrottled = new Date().getTime() - this.lastComponentUpdate > throttleDuration

    if (isNewEvents || wasThrottled) {
      return this.isThrottled
    }

    return false
  }
}

export default EventStream
