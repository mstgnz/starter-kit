import { Component, OnInit } from '@angular/core';
//import { FormControl, FormGroup, Validators } from '@angular/forms';
//import { Router } from '@angular/router';
//import { MessageService } from 'primeng/api';
//import { CustomValidator } from 'src/app/helpers/custom.validator';
//import { ApiService } from 'src/app/services/api.service';

@Component({
  selector: 'app-new-password',
  standalone: false,
  templateUrl: './new-password.component.html',
  styles: ``
})
export class NewPasswordComponent implements OnInit {

  /* public formGroup = new FormGroup(
    {
      email: new FormControl(null, [Validators.required, CustomValidator.isValidEmail]),
      forgot_code: new FormControl(null, [Validators.required]),
      password: new FormControl(null, [Validators.required, Validators.minLength(8)]),
      re_password: new FormControl(null, [Validators.required, Validators.minLength(8)])
    },
    CustomValidator.isPasswordMatch('password', 're_password')
  ) */

  constructor(
    /*  private router: Router,
     private apiService: ApiService,
     private messageService: MessageService */
  ) { }

  ngOnInit(): void { }

  submit() {
    /* if (this.formGroup.valid) {
      this.apiService.forgotPasswordChange(
        this.formGroup.value.email,
        this.formGroup.value.forgot_code,
        this.formGroup.value.password,
        this.formGroup.value.re_password
      ).subscribe(result => {
        if (result.status) {
          this.router.navigateByUrl('login')
          this.messageService.add({ severity: 'success', summary: 'New Password', detail: 'Şifreniz Değiştirildi.', life: 5000 })
        } else {
          this.messageService.add({ severity: 'error', summary: 'New Password', detail: result.message })
        }
      })
    } else {
      this.messageService.add({ severity: 'error', summary: 'New Password', detail: 'Hatalı Form' })
    } */
  }

}
