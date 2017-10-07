import React from 'react';

export default ({
                  file,
                  name,
                }) => {
  return (
    <div>
      <h4>{name}: file {file.default ?
        <span>(default = {file.default})</span> : null}</h4>
      <h5>{file.description}</h5>
    </div>
  );
}
