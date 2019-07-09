import React, { Component } from 'react';

class HelloWorld extends Component {
  constructor(props) {
    super(props);
    this.state = {
      error: null,
      isLoaded: false,
      text: ""
    };
  }

  componentDidMount() {
    fetch("http://localhost:3000/api")
      .then(res => res.text())
      .then(
        (result) => {
          this.setState({
            isLoaded: true,
            text: result
          });
        },
        // Note: it's important to handle errors here
        // instead of a catch() block so that we don't swallow
        // exceptions from actual bugs in components.
        (error) => {
          this.setState({
            isLoaded: true,
            error
          });
        }
      )
  }
  render() {
    const { error, isLoaded, text } = this.state;
    if (error) {
      return <div>Error: {error.message}</div>;
    } else if (!isLoaded) {
      return <div>Loading...</div>;
    } else {
      return (
        <div className="hello-world">
          <h3>{text}</h3>
        </div>
      );
    }
  }
}

export { HelloWorld };
