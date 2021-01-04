import React, { useState } from 'react'
import ReactAutosuggest from 'react-autosuggest'
import brandColors from '../../brandColors'
import { css, cx } from '@emotion/css'
import formElementStyles from '../formElementStyles'

interface Props<TOption> {
  autoFocus?: boolean
  className?: string
  getValue?: (option: TOption) => string
  options: TOption[]
  onSearch: (searchText: string) => any
  onSelect: (data: TOption) => any
  placeholder: string
  render
  renderIcon?: () => any,
  shouldSearch?: (searchText: string) => boolean
}

/**
 * An auto complete input
 */
export default function <TOption>(
  {
    autoFocus,
    className,
    getValue,
    options,
    onSearch,
    onSelect,
    placeholder,
    render,
    renderIcon,
    shouldSearch
  }: Props<TOption>
) {
  const [value, setValue] = useState()

  return (
    <ReactAutosuggest
      focusInputOnSuggestionClick={false}
      getSuggestionValue={
        getValue
          ? getValue
          : option => option.name
      }
      highlightFirstSuggestion
      inputProps={{
        onChange: (e, { newValue }) => setValue(newValue),
        placeholder,
        autoFocus,
        value: value || ''
      }}
      onSuggestionsFetchRequested={({ value }) => onSearch(value)}
      onSuggestionsClearRequested={() => { }}
      onSuggestionSelected={(e, data) => onSelect(data.suggestion)}
      renderSuggestion={render}
      renderInputComponent={inputProps =>
        <div
          className={css({
            display: 'flex',
            alignItems: 'center',
            border: `solid .1rem ${brandColors.lightGray} !important`,
            borderRadius: '.5rem',
            ...options.length
              ? {
                borderBottomLeftRadius: 0,
                borderBottomRightRadius: 0
              }
              : {},
            ...renderIcon
              ? {
                paddingLeft: '1.2rem'
              }
              : {}
          })}
        >
          {
            renderIcon
              ? renderIcon()
              : null
          }
          <input
            {...inputProps}
          />
        </div>
      }
      suggestions={options}
      shouldRenderSuggestions={
        shouldSearch
          ? shouldSearch
          : () => true
      }
      theme={{
        container: cx(
          css({
            width: '100%',
            position: 'relative'
          }),
          className
        ),
        input: css({
          ...formElementStyles,
          backgroundColor: brandColors.white,
          color: brandColors.black,
          padding: '1rem 1.2rem',
          width: '100%'
        }),
        suggestionsContainerOpen: css({
          backgroundColor: brandColors.white,
          border: `1px solid ${brandColors.lightGray}`,
          borderBottomLeftRadius: '.3rem',
          borderBottomRightRadius: '.3rem',
          display: 'block',
          padding: '1rem 1.2rem',
          position: 'absolute',
          left: 0,
          right: 0,
          zIndex: 2
        }),
        suggestionsList: css({
          textAlign: 'left',
          listStyleType: 'none',
          margin: 0,
          padding: 0
        }),
        suggestion: css({
          cursor: 'pointer',
          margin: '.25rem 0',
          ':hover': {
            backgroundColor: brandColors.reallyLightGray
          },
          ':first-child': {
            marginTop: 0
          },
          ':last-child': {
            marginBottom: 0
          }
        })
      }}
    />
  )
}