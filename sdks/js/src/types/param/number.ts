/**
 * Number typed parameter
 */
export default interface Number {
    constraints?: any
    default?: number | string | null | undefined
    description?: string | null | undefined
    isSecret?: boolean | null | undefined
}