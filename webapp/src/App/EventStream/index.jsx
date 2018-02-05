import React, {PureComponent} from 'react';
import Event from './Event';
import eventStore from '../../core/eventStore'
import {toast} from 'react-toastify';
import './index.css';

class EventStream extends PureComponent {
  constructor(props) {
    super(props);

    this.state = {
      events: [],
    };

    this.isScrolledToEnd = true;
    this.eventStreamCloser = () => {
    };
  }

  componentDidMount() {
    this.eventStreamCloser = eventStore.getStream({
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
      onClose: closeEvent => {
        toast.error(
          `event stream closed unexpectedly; msg was ${JSON.stringify(closeEvent)}`
        )
      }
    });
  }

  componentWillUnmount() {
    this.eventStreamCloser();
  }

  componentDidUpdate() {
    if (this.isAtBottom) {
      // only scroll if at bottom; otherwise we'd be overriding the users custom scroll
      this.ref.scrollTop = this.ref.scrollHeight
    }
  }

  componentWillUpdate() {
    this.isAtBottom = (this.ref.scrollHeight - (this.ref.scrollTop + this.ref.offsetHeight)) < 100;
  }

  setRef = ref => this.ref = ref;

  render() {
    // @TODO: shore up performance when many events exist
    return <div className='events' ref={this.setRef}>
      {this.state.events.map((event, index) => <Event key={index} event={event}/>)}
    </div>
  }

}

export default EventStream;
