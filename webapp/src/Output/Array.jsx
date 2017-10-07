import React from 'react';

export default ({
                  array,
                  name,
                }) => {
  return (
    <div>
      <h4>{name}: array {array.default ?
        <span>(default = {JSON.stringify(array.default, null, '\t')})</span> : null}</h4>
      <h5>{array.description}</h5>
    </div>
  );
}
