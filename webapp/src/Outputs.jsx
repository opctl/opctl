import React from 'react';
import Param from './Param';

export default function Outputs(props) {

    const outputs = [];
    if (props.value) {
        Object.entries(props.value).forEach(([name, param]) => {
            outputs.push(<Param name={name} param={param} key={name}/>)
        });
    }

    return (
        <div>
            <h2>Outputs:</h2>
            {outputs}
        </div>
    );
}
