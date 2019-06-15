/**
 * Typed value
 */
export default interface Value {
    /**
     * Array typed value
     */
    array?: any[] | null | undefined

    /**
     * Boolean typed value
     */
    boolean?: boolean | null | undefined

    /**
     * String typed value
     */
    dir?: string | null | undefined

    /**
     * File typed value
     */
    file?: string | null | undefined

    /**
     * Number typed value
     */
    number?: number | null | undefined

    /**
     * Object typed value
     */
    object?: object | null | undefined

    /**
     * Socket typed value
     */
    socket?: string | null | undefined

    /**
     * String typed value
     */
    string?: string | null | undefined
}