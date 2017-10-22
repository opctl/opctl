import React from 'react';
import Markdown from 'react-remarkable';
import highlightJs from 'highlightjs';
import 'highlightjs/styles/github.css';

export default ({value}) => <div className='custom-control-description'>
  <Markdown options={{
    highlight: (str, lang) => {
      if (lang && highlightJs.getLanguage(lang)) {
        try {
          return highlightJs.highlight(lang, str).value;
        } catch (err) {
        }
      }

      try {
        return highlightJs.highlightAuto(str).value;
      } catch (err) {
      }

      return ''; // use external default escaping
    }
  }}>{value}</Markdown>
</div>;
