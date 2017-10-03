import React, {Component} from 'react';
import ReactList from 'react-list'
import Event from './Event';

class EventBrowser extends Component {
  constructor(props) {
    super(props);

    this.state = {
      events: [],
    };
  }

  componentDidMount() {
    this.ws = new WebSocket('ws://localhost:42224/events/stream');
    this.ws.onmessage = msg => {
      this.setState(prevState => ({
        events: [...prevState.events, JSON.parse(msg.data)]
      }));
    };
  }

  componentWillUnmount() {
    this.ws.close()
  }

  renderItem(index, key) {
    const event = this.state.events[index];
    return <Event key={key} event={event}/>;
  }

  render() {
    return (
      <div>
        <h2>Events</h2>
        <div style={{overflow: 'auto', maxHeight: 900}}>
          <ReactList
            itemRenderer={(index, key) => this.renderItem(index, key)}
            length={this.state.events.length}
            type='uniform'
          />
        </div>
      </div>
    );
  }

}

export default EventBrowser;
