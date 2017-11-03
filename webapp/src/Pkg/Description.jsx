import React from 'react';
import Markdown from 'react-remarkable';
import highlightJs from 'highlightjs';
import 'highlightjs/styles/github.css';

export default ({value}) =>
  <Markdown options={{
    highlight: (str, lang) => {
      if (lang && highlightJs.getLanguage(lang)) {
        return highlightJs.highlight(lang, str).value;
      }
      return highlightJs.highlight(lang, str).value;
    }
  }}>
    {value}
  </Markdown>;
