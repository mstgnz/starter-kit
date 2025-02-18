export const isNumeric = (value: string) => {
    return !isNaN(parseInt(value, 10)) && /^\d+$/.test(value)
}

export const isNumeric2 = (n: any) => {
    return !isNaN(parseFloat(n)) && isFinite(n)
}

export const turkishToLower = (value: string) => {
    var letters: { [key: string]: string } = { "İ": "i", "I": "ı", "Ş": "ş", "Ğ": "ğ", "Ü": "ü", "Ö": "ö", "Ç": "ç" };
    value = value.replace(/(([İIŞĞÜÇÖ]))/g, function (letter) { return letters[letter]; })
    return value.toLowerCase();
};

export const turkishToUpper = (value: string) => {
    var letters: { [key: string]: string } = { "i": "İ", "ş": "Ş", "ğ": "Ğ", "ü": "Ü", "ö": "Ö", "ç": "Ç", "ı": "I" };
    value = value.replace(/(([iışğüçö]))/g, function (letter) { return letters[letter]; })
    return value.toUpperCase();
};

export const replaceNullsWithZero = (data: any) => {
    const result = {};
    for (const key in data) {
        if (data.hasOwnProperty(key)) {
            (result as { [key: string]: any })[key] = data[key] === null ? 0 : data[key];
        }
    }
    return result;
}

export const removeEmptyObjects = (obj: { [key: string]: any }) => {
    Object.keys(obj).forEach(key => {
        if (typeof obj[key] === 'object' && obj[key] !== null) {
            removeEmptyObjects(obj[key])
            if (Object.keys(obj[key]).length === 0) {
                delete obj[key]
            }
        }
    })
}