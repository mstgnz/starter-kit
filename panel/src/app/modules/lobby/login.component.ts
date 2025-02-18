import { Component, OnInit } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';

import { Router } from '@angular/router';
import { ApiService } from '../../services/api.service';
import { AuthService } from '../../services/auth.service';
import { AlertifyService } from '../../services/alertify.service';
import { User } from '../../graphql/users.graphql';

@Component({
  selector: 'app-login',
  standalone: false,
  templateUrl: './login.component.html',
  styles: ``
})
export class LoginComponent implements OnInit {

  public loading: boolean = false
  public isCode: boolean = false
  private tryCount: number = 1
  private token: string = ''
  private user: User = {} as User

  public chooseList = [
    { "id": "sms", "name": "Doğrulama Kodunu Sms Gönder" },
    { "id": "email", "name": "Doğrulama Kodunu Email Gönder" }
  ]
  public code = new FormControl<string>('', [Validators.required])
  public formGroup = new FormGroup({
    send: new FormControl('email', [Validators.required]),
    email_or_phone: new FormControl('', [Validators.required]),
    password: new FormControl('', [Validators.required])
  })

  constructor(
    private router: Router,
    private apiService: ApiService,
    private authService: AuthService,
    private alertifyService: AlertifyService
  ) { }

  ngOnInit(): void {
    this.authService.verify().subscribe(response => {
      if (response.success) {
        this.authService.login(response.data.user)
        this.redirect()
      }
    })
  }

  formSubmit() {
    if (this.formGroup.valid) {
      this.loading = true
      const email_or_phone = this.formGroup.value.email_or_phone?.trim() ?? ''
      const password = this.formGroup.value.password?.trim() ?? ''
      const sendMethod = this.formGroup.value.send ?? 'email'
      this.authService.apiLogin(email_or_phone, password, sendMethod).subscribe(response => {
        this.loading = false
        if (response.success) {
          this.isCode = true
          this.token = response.data.token
          this.user = response.data.user
          this.alertifyService.success(response.message)
        } else {
          this.alertifyService.error(response.message)
        }
      })
    }
  }

  checkCode() {
    if (this.tryCount == 3) {
      this.isCode = false
      this.tryCount = 1
      this.alertifyService.error('Doğrulama Kodunu 3 Kere Yanlış Girdiniz!')
      this.router.navigateByUrl('/')
    } else {
      this.loading = true
      const verifyCode = Number(this.code.value ?? '');
      const email_or_phone = this.formGroup.value.email_or_phone ?? '';
      this.apiService.verifyCode(email_or_phone, verifyCode).subscribe(response => {
        this.loading = false
        if (response.success) {
          localStorage.setItem('access_token', this.token)
          this.authService.login(this.user)
          this.redirect()
        } else {
          this.code.patchValue('')
          this.alertifyService.error(response.message)
          this.tryCount++
        }
      })
    }
  }

  redirect() {
    if (this.authService.currentUser.user_type_id == 1) {
      this.router.navigateByUrl('/admin')
    }
    if (this.authService.currentUser.user_type_id === 2) {
      this.router.navigateByUrl('/facility')
    }
  }

}