import React from 'react';

export default ({
                  number,
                  name,
                }) => {
  return (
    <div>
      <h4>{name}: number {number.default ?
        <span>(default = {number.default})</span> : null}</h4>
      <h5>{number.description}</h5>
    </div>
  );
}
