import React from 'react'
import Container from '../Container'
import Fs from '../Fs'
import brandColors from '../../brandColors'

export default () => {
    return (
        <Container
            style={{
                display: 'flex',
                flexDirection: 'column',
                height: '100%'
            }}
        >
            <div
                style={{
                    flexShrink: 0,
                    borderBottom: `solid .1rem ${brandColors.reallyLightGray}`,
                }}
            >
                <h3>Explorer</h3>
            </div>
            <Fs
                style={{
                    borderBottom: `solid .1rem ${brandColors.reallyLightGray}`
                }}
            />
        </Container>
    )
}