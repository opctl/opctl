import React, {Component} from 'react';
import ReactList from 'react-list';
import {getRenderedHeight} from 'react-rendered-size';
import Event from './Event';
import {toast} from 'react-toastify';

const throttleDuration = 200;

class EventBrowser extends Component {
  constructor(props) {
    super(props);

    this.state = {
      events: [],
    };

    this.lastComponentUpdate = new Date().getTime();
  }

  componentDidMount() {
    let queryParts = [];
    if (this.props.filter && this.props.filter.root) {
      queryParts.push(`roots=${encodeURIComponent(this.props.filter.root)}`);
    }

    // @TODO: move to opspecNodeApiClient
    // @TODO: don't assume local node
    this.ws = new WebSocket(`ws://localhost:42224/events/stream?${queryParts.join('&')}`);
    this.ws.onmessage = msg => {
      const event = JSON.parse(msg.data);
      // cache rendered height
      event.height = getRenderedHeight(<Event event={event}/>);

      this.setState(prevState => ({
        events: [...prevState.events, event]
      }));
    };

    this.ws.onerror = error => toast.error(
      `encountered error streaming events; error was ${error.message}`
    );

    // maintain an update interval so throttled renders get re-processed every throttleDuration
    this.interval = setInterval(() => {
      this.setState(prevState => ({
        throttleInterval: new Date().getTime()
      }));
    }, throttleDuration);
  }

  componentDidUpdate() {
    if (this.isAtBottom) {
      // only scroll if at bottom; otherwise we'd be overriding the users custom scroll
      this.messageList.scrollTop = this.messageList.scrollHeight
    }
    this.lastComponentUpdate = new Date().getTime();
  }

  componentWillUnmount() {
    this.ws.close();
    clearInterval(this.interval);
  }

  componentWillUpdate() {
    this.isAtBottom = (this.messageList.scrollHeight - (this.messageList.scrollTop + this.messageList.offsetHeight)) < 100;
  }

  renderItem(index, key) {
    return (<Event key={key} event={this.state.events[index]}/>);
  }

  render() {
    return (
      <div className='events'>
        <div
          style={{overflow: 'auto', height: '100vh'}}
          ref={(div) => {
            this.messageList = div
          }}
        >
          <ReactList
            itemSizeGetter={(index) => this.state.events[index].height}
            itemRenderer={(index, key) => this.renderItem(index, key)}
            length={this.state.events.length}
            type='variable'
          />
        </div>
      </div>
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

export default EventBrowser;
