import React, { ReactNode } from 'react'
import ReactModal from 'react-modal'
import { cx, css } from 'emotion'
import LeftNav from './LeftNav'
import RightNav from './RightNav'
import brandColors from '../../brandColors'

interface Props {
  children: any
  className?: string | null
  isCompleteDisabled?: boolean
  isOpen: boolean
  onComplete?: (value: any) => void
  onClose?: () => void
  title?: ReactNode | null
}

// required per (http://reactcommunity.org/react-modal/accessibility/)
// @TODO: re-enable & set ariaHideApp={true}
// ReactModal.setAppElement('#root')
export default (
  {
    children,
    className,
    isCompleteDisabled,
    isOpen,
    onComplete,
    onClose,
    title,
  }: Props
) => <ReactModal
  className={cx(
    css({
      '@media screen and (min-width: 600px)': {
        minWidth: '600px',
        borderRadius: '.3rem'
      },
      '@media screen and (max-width: 600px)': {
        width: '100%',
        height: '100%'
      },
      top: '50%',
      left: '50%',
      right: 'auto',
      bottom: 'auto',
      marginRight: '-50%',
      transform: 'translate(-50%, -50%)',
      position: 'absolute',
      border: '1px solid #ccc',
      background: brandColors.white,
      display: 'flex',
      flexDirection: 'column',
      outline: 'none',
      maxHeight: '100%',
      overflowY: 'auto'
    }),
    className
  )}
  style={{
    overlay: {
      backgroundColor: 'rgba(190, 190, 190, 0.75)',
      zIndex: 1000
    }
  }}
  isOpen={isOpen}
  onRequestClose={onClose}
  ariaHideApp={false}
>
    {
      title || onClose || onComplete
        ? <div
          className={css({
            borderBottom: `solid thin ${brandColors.lightGray}`,
            display: 'flex',
            alignItems: 'center',
            height: '3rem',
            width: '100%',
            flexShrink: 0,
            position: 'static'
          })}
        >
          <div
            className={css({
              color: brandColors.blue,
              cursor: 'pointer',
              flex: '1 0 0',
              marginLeft: '1.2rem'
            })}
          >
            {
              onClose &&
              <div
                onClick={onClose}
              >
                Close
              </div>
            }
          </div>
          <h2
            className={css({
              flex: 'none',
              marginBottom: 0,
              color: brandColors.black
            })}
          >
            {title}
          </h2>
          <RightNav
            isCompleteDisabled={isCompleteDisabled}
            onClick={onComplete}
          />
        </div>
        : null
    }
    {children}
  </ReactModal>