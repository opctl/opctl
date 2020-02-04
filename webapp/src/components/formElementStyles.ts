import { ObjectInterpolation } from "emotion"

const styles: ObjectInterpolation<undefined> = {
  fontSize: '1rem',
  maxWidth: '100vw',
  border: 0,
  borderRadius: '.5rem',
  // remove blue outline (on focus) in chrome
  outline: '0px !important',
  WebkitAppearance: 'none',
  'input::placeholder ::placeholder': {
    // use input color
    color: 'inherit',
    opacity: .5
  },
  ':disabled': {
    opacity: .5
  }
}

export default styles
