import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'percent'
})
export class PercentPipeComponent implements PipeTransform {

    transform(number: number): string {
        number = number < 1 ? number * 100 : number
        return "%" + number.toFixed(2)
    }

}