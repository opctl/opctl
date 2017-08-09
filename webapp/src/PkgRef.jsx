import React from 'react';

export default function PkgRef(props){
    return (
        <span>{props.name}#{props.version}</span>
    )
}