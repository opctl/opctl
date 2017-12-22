import React from 'react';
import Markdown from '../Markdown';
import 'highlightjs/styles/github.css';

export default ({value, pkgRef}) =>
  <div className='custom-control-description'>
    <Markdown value={value} pkgRef={pkgRef}/>
  </div>;
