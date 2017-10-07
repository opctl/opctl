import React from 'react';

export default (props) => {
    return (
        <div>
            <h4>{props.name}: socket</h4>
            <h5>{props.socket.description}</h5>
        </div>
    );
}
