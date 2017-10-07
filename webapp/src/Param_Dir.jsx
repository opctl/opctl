import React from 'react';

export default (props) => {
    return (
        <div>
            <h4>{props.name}: dir {props.dir.default ? <span>(default = {props.dir.default})</span> : null }</h4>
            <h5>{props.dir.description}</h5>
        </div>
    );
}
