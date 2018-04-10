import React from 'react'
import Markdown from 'react-remarkable'
import highlightJs from 'highlightjs'
import 'highlightjs/styles/github.css'

export default ({value, opRef}) => {
  value = !value ? '' : value.replace(
    /]\(\/.+\)/,
    match => {
      const contentPath = match.slice(2, match.length - 1)
      return `](/api/ops/${encodeURIComponent(opRef)}/contents/%2f${encodeURIComponent(contentPath)})`
    })

  return (<Markdown options={{
    highlight: (str, lang) => {
      if (lang && highlightJs.getLanguage(lang)) {
        return highlightJs.highlight(lang, str).value
      }
      return highlightJs.highlightAuto(str).value
    }
  }}>
    {value}
  </Markdown>)
}
