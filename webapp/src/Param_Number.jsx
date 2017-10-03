import React from 'react';

export default function Param(props) {
    return (
        <div>
            <h4>{props.name}: number {props.number.default ? <span>(default = {props.number.default})</span> : null }</h4>
            <h5>{props.number.description}</h5>
        </div>
    );
}
