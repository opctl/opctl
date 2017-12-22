import React, {Component} from 'react';
import ReactList from 'react-list';
import {getRenderedHeight} from 'react-rendered-size';
import Event from './Event';
import opspecNodeApiClient from '../../utils/clients/opspecNodeApi';
import {toast} from 'react-toastify';
import './index.css';

const throttleDuration = 200;

class EventBrowser extends Component {
  constructor(props) {
    super(props);

    this.state = {
      events: [],
    };

    this.eventStreamCloser = () => {};
    this.lastComponentUpdate = new Date().getTime();
  }

  componentDidMount() {
    this.eventStreamCloser = opspecNodeApiClient.event_stream_get({
      filter: this.props.filter,
      onEvent: rawEvent => {
        const event = Object.assign(
          {
            // cache rendered height
            height: getRenderedHeight(<Event event={rawEvent}/>),
          },
          rawEvent
        );

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
    if (this.isAtBottom) {
      // only scroll if at bottom; otherwise we'd be overriding the users custom scroll
      this.messageList.scrollTop = this.messageList.scrollHeight
    }
    this.lastComponentUpdate = new Date().getTime();
  }

  componentWillUnmount() {
    this.eventStreamCloser();
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
          className='event-browser'
          style={{overflow: 'auto'}}
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
