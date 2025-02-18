import { AbstractControl, FormGroup, ValidationErrors } from "@angular/forms";

export class CustomValidator {

    // IMAGE
    static isValidExtension(input: AbstractControl): ValidationErrors | null {
        const value = input.value as string;

        if (value.endsWith('.jpg') || value.endsWith('.jpeg') || value.endsWith('.png')) {
            return null;
        }
        return {
            wrongExtension: true,
            availableExtension: "jpg or png"
        }
    }

    // EMAIL
    static isValidEmail(input: AbstractControl): ValidationErrors | null {
        const value = input.value as string;
        const regex = /^\w+([\.-]?\w+)*@\w+([\.-]?\w+)*(\.\w{2,3})+$/;
        const valid = regex.test(value);

        if (valid) {
            return null;
        }
        return {
            wrongEmail: true
        }
    }

    // FIELD MATCH
    static isPasswordMatch(controlName: string, matchingControlName: string): ValidationErrors | null {
        return (formGroup: FormGroup) => {
            const control = formGroup.controls[controlName]
            const matchingControl = formGroup.controls[matchingControlName]
            if (matchingControl.errors) {
                return
            }
            // set error on matchingControl if validation fails
            if (control.value !== matchingControl.value) {
                matchingControl.setErrors({ mustMatch: true })
            } else {
                matchingControl.setErrors(null)
            }
            return null
        }
    }

    // DATE GRATHER THEN NOW
    static isGratherThanNow(input: AbstractControl): ValidationErrors | null {
        const value = input.value as Date
        if (value && new Date(value) >= new Date()) {
            return null
        }
        return {
            wrongDate: true
        }
    }

    // DATE GRATHER THEN DATE
    static isGratherThanDate(startDate: string, finishDate: string): ValidationErrors | null {
        return (formGroup: FormGroup) => {
            const start = formGroup.controls[startDate]
            const finish = formGroup.controls[finishDate]
            if (start.errors || finish.errors) {
                return null
            }
            if (start.value >= finish.value) {
                return { wrongDate: true }
            }
            return null
        }
    }

}