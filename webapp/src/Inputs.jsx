import React from 'react';
import Input from './Input';
import opSpecNodeApiClient from './opspecNodeApiClient'
import {toast} from 'react-toastify';

export default ({
                  value,
                  pkgRef,
                }) => {

  const args = {};
  const inputs = [];
  if (value) {
    Object.entries(value).forEach(([name, input]) => {
      inputs.push(<Input
        onChange={value => args[name] = value}
        name={name}
        input={input}
        key={name}
      />);
    });
  }

  const handleSubmit = (e) => {
    e.preventDefault();

    const req = {
      args,
      pkg: {
        ref: pkgRef,
      }
    };

    opSpecNodeApiClient.startOp(req)
      .catch(error => {
        toast.error(error.message);
      });
  };

  return (
    <div>
      <h2>Inputs:</h2>
      <form onSubmit={handleSubmit}>
        {inputs}
        <input className='btn btn-primary btn-lg' id='startOp_Submit' type='submit' value='start'/>
      </form>
    </div>
  );
}
