import React from 'react';

export default function Param(props) {
    return (
        <div>
            <h4>{props.name}: file {props.file.default ? <span>(default = {props.file.default})</span> : null }</h4>
            <h5>{props.file.description}</h5>
        </div>
    );
}
