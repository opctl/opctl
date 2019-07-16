import objectUnderTest from './get'
import { Duplex } from 'stream'

const nock = require('nock')
const providedNetRef = 'http://dummyNetRef'

afterEach(() => {
    nock.cleanAll()
})

describe('dataGet', () => {
    it('makes expected http request', async () => {
        /* arrange */
        const providedDataRef = '//dummyDataRef'

        const scope = nock(providedNetRef)
            .log(console.log)
            .get(
                `/data/${encodeURIComponent(providedDataRef)}`)
            .reply('200')

        /* act */
        await objectUnderTest(
            providedNetRef,
            providedDataRef
        )

        /* assert */
        scope.isDone()
    })
    describe('response is 300', () => {
        it('returns expected result', async () => {
            /* arrange */
            const providedDataRef = '//dummyDataRef'
            const expectedErrMsg = 'dummyErrorMsg'

            nock(providedNetRef)
                .log(console.log)
                .get(
                    `/data/${encodeURIComponent(providedDataRef)}`)
                .reply('300', expectedErrMsg)

            /* act */
            let actualErr
            try {
                await objectUnderTest(
                    providedNetRef,
                    providedDataRef
                )
            } catch (err) {
                actualErr = err
            }

            /* assert */
            expect(actualErr.message).toEqual(expectedErrMsg)
        })
    })
    describe('response is 200', () => {
        it('returns expected result', async () => {
            /* arrange */
            const providedDataRef = '//dummyDataRef'

            const expectedContent = 'dummyContent'

            nock(providedNetRef)
                .log(console.log)
                .get(
                    `/data/${encodeURIComponent(providedDataRef)}`, () => {
                        let stream = new Duplex()
                        stream.push(Buffer.from(expectedContent))
                        stream.push(null)
                        return stream
                    })
                .reply('200', expectedContent)

            /* act */
            const actualContent = await objectUnderTest(
                providedNetRef,
                providedDataRef
            ).then(response => response.text())

            /* assert */
            expect(actualContent).toEqual(expectedContent)
        })
    })
})
