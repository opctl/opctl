import React from 'react';
import Param from './Param';

export default function Outputs(props) {

    const outputs = [];
    if (props.value) {
        Object.entries(props.value).forEach(([name, type]) => {
            outputs.push(<Param name={name} type={type} key={name}/>)
        });
    }

    return (
        <div>
            <h2>Outputs:</h2>
            {outputs}
        </div>
    );
}