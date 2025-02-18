export const _unique = <T>(array: T[]) => [...new Set(array)]

export const _intersection = <T>(array1: T[], array2: T[]) => array1.filter((v) => array2.includes(v))

export const _diff = <T>(array1: T[], array2: T[]) => array1.filter((v) => !array2.includes(v))

export const symDiff = <T>(array1: T[], array2: T[]) => _diff(array1, array2).concat(_diff(array2, array1))

export const union = <T>(array1: T[], array2: T[]) => _diff(array1, array2).concat(array2)

export const _intersectionBy = <T>(
    array1: T[],
    array2: T[],
    predicate: (array1Value: T, array2Value: T) => boolean
) => array1.filter((v) => array2.some((u) => predicate(v, u)))

export const _diffBy = <T>(
    array1: T[],
    array2: T[],
    predicate: (array1Value: T, array2Value: T) => boolean
) => array1.filter((v) => !array2.some((u) => predicate(v, u)))

export const _uniqueBy = <T>(
    array: T[],
    predicate: (v: T, i: number, a: T[]) => string
) =>
    Object.values(
        array.reduce((acc, value, index) => {
            acc[predicate(value, index, array)] = value
            return acc
        }, {} as { [key: string]: T })
    )
