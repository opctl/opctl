import Array from './array'
import Boolean from './boolean'
import Dir from './dir'
import File from './file'
import Number from './number'
import Object from './object'
import Socket from './socket'
import String from './string'

/**
 * Typed parameter of an op
 */
export default interface Param {
    array?: Array | null | undefined
    boolean?: Boolean | null | undefined
    dir?: Dir | null | undefined
    file?: File | null | undefined
    number?: Number | null | undefined
    object?: Object | null | undefined
    socket?: Socket | null | undefined
    string?: String | null | undefined
}
