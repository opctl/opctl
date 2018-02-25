import React from 'react';

export default ({value, pkgRef}) => {
  if (value) {
    value = value.replace(
      /^\/.+$/,
      match => {
        const contentPath = match.slice(0, match.length);
        return `/api/pkgs/${encodeURIComponent(pkgRef)}/contents/%2f${encodeURIComponent(contentPath)}`
      });
  } else {
    value = "opspec-icon.svg"
  }

  return (<img src={value} alt={'icon'} style={{height: '10vw'}} />);
}
