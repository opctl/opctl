import React from 'react';

export default ({
                  dir,
                  name,
                }) => {
  return (
    <div>
      <h4>{name}: dir {dir.default ?
        <span>(default = {dir.default})</span> : null}</h4>
      <h5>{dir.description}</h5>
    </div>
  );
}
