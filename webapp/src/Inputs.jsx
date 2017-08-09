import React from 'react';
import Param from './Param';

export default function Inputs(props) {

    const inputs = [];
    if (props.value) {
        Object.entries(props.value).forEach(([name, type]) => {
            inputs.push(<Param name={name} type={type} key={name}/>)
        });
    }

    return (
        <div>
            <h2>Inputs:</h2>
            {inputs}
        </div>
    );
}