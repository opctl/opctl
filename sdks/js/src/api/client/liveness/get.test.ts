import objectUnderTest from './get'

const nock = require('nock')
const providedNetRef = 'http://dummyNetRef'

afterEach(() => {
    nock.cleanAll()
})

describe('livenessGet', () => {
    it('makes expected http request', async () => {
        /* arrange */
        const scope = nock(providedNetRef)
            .log(console.log)
            .get(
                `/liveness`)
            .reply('200')

        /* act */
        await objectUnderTest(
            providedNetRef
        )

        /* assert */
        scope.isDone()
    })
    describe('response is 300', () => {
        it('returns expected result', async () => {
            /* arrange */
            const expectedErrMsg = 'dummyErrorMsg'

            nock(providedNetRef)
                .log(console.log)
                .get(
                    `/liveness`)
                .reply('300', expectedErrMsg)

            /* act */
            let actualErr: null | Error = null
            try {
                await objectUnderTest(
                    providedNetRef
                )
            } catch (err) {
                actualErr = err as Error
            }

            /* assert */
            expect(actualErr?.message).toEqual(expectedErrMsg)
        })
    })
    describe('response is 200', () => {
        it('returns expected result', async () => {
            /* arrange */
            nock(providedNetRef)
                .log(console.log)
                .get(
                    `/liveness`)
                .reply('200')

            /* act */
            const actualResponse = await objectUnderTest(
                providedNetRef
            )

            /* assert */
            expect(actualResponse).toEqual('')
        })
    })
})