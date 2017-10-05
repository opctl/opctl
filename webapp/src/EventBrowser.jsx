import React, {Component} from 'react';
import ReactList from 'react-list';
import {getRenderedHeight} from 'react-rendered-size';
import Event from './Event';

class EventBrowser extends Component {
  constructor(props) {
    super(props);

    this.state = {
      events: [],
    };
  }

  componentDidMount() {
    // @TODO: don't assume local node
    this.ws = new WebSocket('ws://localhost:42224/events/stream');
    this.ws.onmessage = msg => {
      const event = JSON.parse(msg.data);
      // cache rendered height
      event.height = getRenderedHeight(<Event event={event}/>);

      this.setState(prevState => ({
        events: [...prevState.events, event]
      }));
    };

    this.ws.onerror = err => {
      console.log(`websocket erred; err was ${err}`);
    };

    this.ws.onclose = ({code, reason}) => {
      console.log(`websocket closed unexpectedly. code, reason were:  ${code}, ${reason}`)
    }
  }

  componentWillUnmount() {
    this.ws.close()
  }

  renderItem(index, key) {
    return (<Event key={key} event={this.state.events[index]}/>);
  }

  render() {
    return (
      <div className='events'>
        <div style={{overflow: 'auto', height: '100vh'}} ref={(div) => {this.messageList = div}}>
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

  componentWillUpdate(){
    this.isAtBottom = (this.messageList.scrollHeight - (this.messageList.scrollTop + this.messageList.offsetHeight)) < 100;
  }

  componentDidUpdate() {
    if (this.isAtBottom) {
      // only scroll if at bottom; otherwise we'd be overriding the users custom scroll
      this.messageList.scrollTop = this.messageList.scrollHeight
    }
  }

}

export default EventBrowser;
