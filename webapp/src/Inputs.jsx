import React from 'react';
import Input from './Input';

export default ({value, onArgChange}) => {
  const inputs = Object.entries(value || {}).map(([name, input]) => (
    <Input
      onArgChange={value => (onArgChange(name, value))}
      name={name}
      input={input}
      key={name}
    />
  ));

  return (
    <div>
      <h2>Inputs:</h2>
      {inputs}
    </div>
  );
}
