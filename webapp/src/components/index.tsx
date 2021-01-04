import React from 'react'
import Workspace from './Workspace'
import { BrowserRouter } from 'react-router-dom'
import { ToastContainer } from 'react-toastify'

// shared vendor css
import 'react-toastify/dist/ReactToastify.min.css'
import 'reboot.css'

export default function components() {
  return (
    <BrowserRouter>
      <ToastContainer
        autoClose={20000}
        style={{ zIndex: 100000 }}
      />
      <Workspace />
    </BrowserRouter>
  )
}