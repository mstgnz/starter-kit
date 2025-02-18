import { Pipe, PipeTransform } from '@angular/core';

@Pipe({
    name: 'string'
})
export class StringPipeComponent implements PipeTransform {

    transform(number: Number): string {
        return number.toString()
    }

}