import React from 'react';

export default (props) => {
    return (
        <div>
            <h4>{props.name}: array {props.array.default ? <span>(default = {JSON.stringify(props.array.default, null, '\t')})</span> : null }</h4>
            <h5>{props.array.description}</h5>
        </div>
    );
}
