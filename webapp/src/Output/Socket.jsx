import React from 'react';

export default ({
                  socket,
                  name,
                }) => {
  return (
    <div>
      <h4>{name}: socket {socket.default ?
        <span>(default = {socket.default})</span> : null}</h4>
      <h5>{socket.description}</h5>
    </div>
  );
}
