/**
 * Object typed parameter
 */
export default interface ParamObject {
    constraints?: any
    default?: object | string | null | undefined
    description?: string | null | undefined
    isSecret?: boolean | null | undefined
}