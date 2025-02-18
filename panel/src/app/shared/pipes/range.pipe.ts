import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'range'
})
export class RangePipeComponent implements PipeTransform {

    transform(length: number, offset: number = 0): number[] {
        if (!length) {
            return [];
        }
        const array: number[] = [];
        for (let n = 1; n <= length; ++n) {
            array.push(offset + n);
        }
        return array;
    }

}