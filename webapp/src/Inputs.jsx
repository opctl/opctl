import React from 'react';
import Param from './Param';

export default function Inputs(props) {

    const inputs = [];
    if (props.value) {
        Object.entries(props.value).forEach(([name, param]) => {
            inputs.push(<Param name={name} param={param} key={name}/>)
        });
    }

    return (
        <div>
            <h2>Inputs:</h2>
            {inputs}
        </div>
    );
}
