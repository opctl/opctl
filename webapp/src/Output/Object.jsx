import React from 'react';

export default ({
                  name,
                  object,
                }) => {
  return (
    <div>
      <h4>{name}: object {object.default ?
        <span>(default = {JSON.stringify(object.default, null, '\t')})</span> : null}</h4>
      <h5>{object.description}</h5>
    </div>
  );
}
