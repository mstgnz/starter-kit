import { turkishToLower } from "./utils"

export const filterName = (lists: any[], name: string): any[] => {
    const filter: any[] = []
    lists.forEach(list => {
        if (turkishToLower(list.name.toString()).indexOf(turkishToLower(name)) == 0) {
            filter.push(list)
        }
    })
    return filter
}