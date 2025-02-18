export class DateHelper {

    // get date with format 2022/11/23
    static getToday(): string {
        return new Date().toISOString().split('T')[0]
    }

    // get current week first date
    static getWeekFirstDate(date: string | undefined = undefined): string {
        let today = date ? new Date(date) : new Date()
        today = new Date(today.getFullYear(), today.getMonth(), today.getDate() - today.getDay() + 1)
        return `${today.getFullYear()}-${today.getMonth() + 1}-${today.getDate()}`
    }

    // get current week last date
    static getWeekLastDate(date: string | undefined = undefined): string {
        let today = date ? new Date(date) : new Date()
        today.setDate(today.getDate() + 7)
        return `${today.getFullYear()}-${today.getMonth() + 1}-${today.getDate()}`
    }

    // get current month last day
    static getMonthFirstDate(date: string | undefined = undefined): string {
        let today = date ? new Date(date) : new Date()
        let month = today.getMonth() + 1
        return `${today.getFullYear()}-${month < 10 ? "0" + month : month}-01`
    }

    // get current month last day
    static getMonthLastDate(date: string | undefined = undefined): string {
        let today = new Date(date ? new Date(date) : new Date())
        today = new Date(today.getFullYear(), today.getMonth() + 1, 0)
        let month = today.getMonth() + 1
        return `${today.getFullYear()}-${month < 10 ? "0" + month : month}-${today.getDate()}`
    }

    static getNextYear(date: string): string {
        let newDate = new Date(date)
        newDate.setFullYear(newDate.getFullYear() + 1)
        return `${newDate.getFullYear()}-${this.fixMonthAndDay(newDate.getMonth() + 1)}-${this.fixMonthAndDay(newDate.getDate())}`
    }

    static fixMonthAndDay(monthOrDate: number): string | number {
        return monthOrDate < 10 ? "0" + monthOrDate : monthOrDate
    }

    static addMinuteDate(date: string, minute: number): Date {
        let newDate = new Date(date)
        newDate.setMinutes(newDate.getMinutes() + minute)
        return newDate
    }

    static dateFormat(date: string): string {
        const newDate = new Date(date)
        return newDate.toString() !== "Invalid Date" ? `${newDate.getFullYear()}-${this.fixMonthAndDay(newDate.getMonth() + 1)}-${this.fixMonthAndDay(newDate.getDate())}` : date
    }

    static diffDateMinute(beforeDate: any, afterDate: any): number {
        beforeDate = new Date(beforeDate)
        afterDate = new Date(afterDate)
        const diffMiliSecond = afterDate - beforeDate
        return diffMiliSecond / (1000 * 60)
    }

    static dateFormatNoSecond(date: string) {
        const newDate = new Date(date)
        return newDate.toString() !== "Invalid Date" ? `${newDate.getFullYear()}-${this.fixMonthAndDay(newDate.getMonth() + 1)}-${this.fixMonthAndDay(newDate.getDate())} ${this.fixMonthAndDay(newDate.getHours())}:${this.fixMonthAndDay(newDate.getMinutes())}` : date
    }

    static dateFormatFull(date: string) {
        const newDate = new Date(date)
        return newDate.toString() !== "Invalid Date" ? `${newDate.getFullYear()}-${this.fixMonthAndDay(newDate.getMonth() + 1)}-${this.fixMonthAndDay(newDate.getDate())} ${this.fixMonthAndDay(newDate.getHours())}:${this.fixMonthAndDay(newDate.getMinutes())}:${this.fixMonthAndDay(newDate.getSeconds())}` : date
    }

}