import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'dateFormat'
})
export class DetePipeComponent implements PipeTransform {

    transform(date: Date, formatOption: string = 'default'): string {
        if (date) {
            date = new Date(date)
            let dateStr = '';
            switch (formatOption) {
                case 'noHour':
                    dateStr = this.getDay(date) + "/" + this.getMonth(date) + "/" + this.getFullYear(date);
                    break;
                case 'noMinute':
                    dateStr = this.getDay(date) + "/" + this.getMonth(date) + "/" + this.getFullYear(date) + " " + this.getHours(date) + ":" + '00';
                    break;
                case 'noSecond':
                    dateStr = this.getDay(date) + "/" + this.getMonth(date) + "/" + this.getFullYear(date) + " " + this.getHours(date) + ":" + this.getMinutes(date);
                    break;
                default:
                    dateStr = this.getDay(date) + "/" + this.getMonth(date) + "/" + this.getFullYear(date) + " " + this.getHours(date) + ":" + this.getMinutes(date) + ":" + this.getSeconds(date);
                    break;
            }
            return dateStr;
        }
        return ""
    }

    getDay(date: Date) {
        let day = date.getDate()
        return day < 10 ? "0" + day : day
    }

    getMonth(date: Date) {
        let month = date.getMonth() + 1
        return month < 10 ? "0" + month : month
    }

    getFullYear(date: Date) {
        return date.getFullYear()
    }

    getHours(date: Date) {
        let hour = date.getHours()
        return hour < 10 ? "0" + hour : hour
    }

    getMinutes(date: Date) {
        let minute = date.getMinutes()
        return minute < 10 ? "0" + minute : minute
    }

    getSeconds(date: Date) {
        let second = date.getSeconds()
        return second < 10 ? "0" + second : second
    }

}