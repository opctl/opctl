import React from 'react'
import ReactRemarkable from 'react-remarkable'
import highlightJs from 'highlightjs'
import 'highlightjs/styles/github.css'

interface Props {
  opRef: string
  value: string | null | undefined
}

export default (
  {
    value,
    opRef
  }: Props
) => {
  value = !value ? '' : value.replace(
    /]\(\/.+\)/,
    match => {
      const contentPath = match.slice(2, match.length - 1)
      return `](/api/ops/${encodeURIComponent(opRef)}/contents/%2f${encodeURIComponent(contentPath)})`
    })

  return (
    <ReactRemarkable
      options={{
        highlight: (
          str: string,
          lang: string
        ) => {
          if (lang && highlightJs.getLanguage(lang)) {
            return highlightJs.highlight(lang, str).value
          }
          return highlightJs.highlightAuto(str).value
        }
      }}
    >
      {value}
    </ReactRemarkable>
  )
}
