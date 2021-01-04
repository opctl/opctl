import React, { Fragment } from 'react'
import HasCall, { CallSerial } from '../HasCall'
import brandColors from '../../brandColors'
import AddCallPopper from '../AddCallPopper'
import { ReactComponent as PlusIcon } from '../../icons/Plus.svg'

interface Props {
  callSerial: CallSerial
  parentOpRef: string
}

export default function CallHasSerial(
  {
    callSerial,
    parentOpRef
  }: Props
) {
  return (
    <div>
      {
        callSerial!.map(
          (childCall, index, array) =>
            <div
              key={index}
              style={{
                minWidth: 0,
                flex: '1 0 0',
                display: 'flex',
                marginLeft: '2rem',
                marginRight: '2rem',
                alignItems: 'center',
                flexDirection: 'column'
              }}
            >
              <HasCall
                key={index}
                call={childCall}
                parentOpRef={parentOpRef}
              />
              {
                index + 1 < array.length
                  ? <Fragment>
                    <div
                      style={{
                        backgroundColor: brandColors.lightGray,
                        height: '100%',
                        minHeight: '2.5rem',
                        width: '.1rem'
                      }}
                    ></div>
                    <AddCallPopper>
                      <PlusIcon
                        style={{
                          backgroundColor: brandColors.white,
                          cursor: 'pointer',
                          fill: brandColors.active,
                          display: 'block'
                        }}
                      />
                    </AddCallPopper>
                    <div
                      style={{
                        backgroundColor: brandColors.lightGray,
                        height: '100%',
                        minHeight: '2.5rem',
                        width: '.1rem'
                      }}
                    ></div>
                  </Fragment>
                  // is leaf
                  : null
              }
            </div>
        )
      }
    </div>
  )
}