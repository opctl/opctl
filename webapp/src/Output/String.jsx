import React from 'react';

export default ({
                  string,
                  name,
                }) => {
  return (
    <div>
      <h4>{name}: string {string.default ?
        <span>(default = {string.default})</span> : null}</h4>
      <h5>{string.description}</h5>
    </div>
  );
}
