import React from 'react';

export default (props) => {
    return (
        <div>
            <h4>{props.name}: string {props.string.default ? <span>(default = {props.string.default})</span> : null }</h4>
            <h5>{props.string.description}</h5>
        </div>
    );
}
