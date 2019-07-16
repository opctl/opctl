/**
 * Array typed parameter
 */
export default interface Array {
    constraints?: any
    default?: any[] | string | null | undefined
    description?: string | null | undefined
    isSecret?: boolean | null | undefined
}