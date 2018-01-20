import React, { Component } from 'react';
import Event from './Event';
import opspecNodeApiClient from '../../utils/clients/opspecNodeApi';
import { AutoSizer, CellMeasurer, CellMeasurerCache, List } from 'react-virtualized';
import { toast } from 'react-toastify';
import './index.css';

const throttleDuration = 200;

class EventStream extends Component {
  constructor(props) {
    super(props);

    this.state = {
      events: [],
    };

    this.mostRecentWidth = 0;
    this.isScrolledToEnd = true;
    this.eventStreamCloser = () => { };
    this.lastComponentUpdate = new Date().getTime();
    this.cache = new CellMeasurerCache({
      fixedWidth: true
    });
  }

  componentDidMount() {
    this.eventStreamCloser = opspecNodeApiClient.event_stream_get({
      filter: this.props.filter,
      onEvent: event => {
        this.setState(prevState => ({
          events: [...prevState.events, event]
        }));
      },
      onError: err => {
        toast.error(
          `encountered error streaming events; error was ${JSON.stringify(err)}`
        )
      },
    });

    // maintain an update interval so throttled renders get re-processed every throttleDuration
    this.interval = setInterval(() => {
      this.setState(() => ({
        throttleInterval: new Date().getTime()
      }));
    }, throttleDuration);
  }

  componentDidUpdate() {
    this.lastComponentUpdate = new Date().getTime();
  }

  componentWillUnmount() {
    this.eventStreamCloser();
    clearInterval(this.interval);
  }

  rowRenderer({ index, isScrolling, key, parent, style }) {
    return (
      <CellMeasurer
        cache={this.cache}
        columnIndex={0}
        key={key}
        parent={parent}
        rowIndex={index}
        width={this.mostRecentWidth}
      >
        <div style={style}>
          <Event event={this.state.events[index]} />
        </div>
      </CellMeasurer>
    );
  }

  handleResize({ width }) {
    if (this.mostRecentWidth !== width) {
      this.mostRecentWidth = width;
      this.cache.clearAll();
      this.list.recomputeRowHeights();
    }
  }

  handleRowsRendered({ stopIndex }) {
    this.isScrolledToEnd = this.state.events.length - stopIndex <= 1;
  }

  render() {
    return (
      this.state.events.length !== 0
        ? <AutoSizer onResize={resized => this.handleResize(resized)}>
          {({ height, width }) => {
            const optionalProps = {};
            if (this.isScrolledToEnd) {
              optionalProps.scrollToIndex = this.state.events.length;
            }
            return (
              <List
                className='events'
                height={height}
                width={width}
                ref={ref => {
                  this.list = ref;
                }}
                onRowsRendered={rowsRendered => this.handleRowsRendered(rowsRendered)}
                deferredMeasurementCache={this.cache}
                rowHeight={this.cache.rowHeight}
                rowRenderer={(renderOpts) => this.rowRenderer(renderOpts)}
                rowCount={this.state.events.length}
                {...optionalProps}
              />
            );
          }}
        </AutoSizer>
        : null
    );
  }

  shouldComponentUpdate(nextProps, nextState) {
    // there are typically so many events received it floods react & results in a completely unresponsive UI
    // we throttle renders to avoid this
    const wasThrottled = this.isThrottled;
    const isNewEvents = this.state.events.length < nextState.events.length;
    this.isThrottled = new Date().getTime() - this.lastComponentUpdate > throttleDuration;

    if (isNewEvents || wasThrottled) {
      return this.isThrottled;
    }

    return false;
  }

}

export default EventStream;
