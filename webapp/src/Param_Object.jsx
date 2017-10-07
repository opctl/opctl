import React from 'react';

export default (props) => {
    return (
        <div>
            <h4>{props.name}: object {props.object.default ? <span>(default = {JSON.stringify(props.object.default, null, '\t')})</span> : null }</h4>
            <h5>{props.object.description}</h5>
        </div>
    );
}
