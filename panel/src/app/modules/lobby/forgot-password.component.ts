import { Component, OnInit } from '@angular/core';
import { FormControl, Validators } from '@angular/forms';
import { CustomValidator } from '../../helpers/custom.validator';
import { Router } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { AlertifyService } from '../../services/alertify.service';

@Component({
  selector: 'app-forgot-password',
  standalone: false,
  templateUrl: './forgot-password.component.html',
  styles: ``
})
export class ForgotPasswordComponent implements OnInit {

  email = new FormControl(null, [Validators.required, CustomValidator.isValidEmail])
  send = new FormControl('email', [Validators.required])

  public chooseList = [
    { "id": "sms", "name": "Doğrulama Kodunu Sms Gönder" },
    { "id": "email", "name": "Doğrulama Kodunu Email Gönder" }
  ]

  constructor(
    private router: Router,
    private apiService: ApiService,
    private alertifyService: AlertifyService
  ) { }

  ngOnInit(): void { }

  submit() {
    if (this.email.valid && this.send.valid) {
      this.apiService.forgotPassword(this.email.value!, this.send.value!).subscribe(result => {
        if (result.status) {
          this.router.navigateByUrl('new-password')
          this.alertifyService.success(result.message)
        } else {
          this.alertifyService.error(result.message)
        }
      })
    } else {
      this.alertifyService.error('Hatalı Form')
    }
  }

}
