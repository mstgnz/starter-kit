import { Pipe, PipeTransform } from "@angular/core";

@Pipe({
    name: "phone"
})
export class PhonePipeComponent implements PipeTransform {

    transform(phone: any) {

        let newPhone = "";

        if (phone) {
            if (phone.length == 11) {
                newPhone += phone.slice(0, 4) + "-" + phone.slice(4, 7) + "-" + phone.slice(7, 9) + "-" + phone.slice(9, 11)
            }

            if (phone.length == 10) {
                newPhone += phone.slice(0, 3) + "-" + phone.slice(3, 6) + "-" + phone.slice(6, 8) + "-" + phone.slice(8, 10)
            }

        }

        return newPhone
    }

}