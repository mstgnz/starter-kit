import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';

import { AuthService } from '../services/auth.service';
import { AlertifyService } from '../services/alertify.service';

@Injectable()
export class LoginGuard implements CanActivate {

  constructor(
    private router: Router,
    private authService: AuthService,
    private alertifyService: AlertifyService
  ) { }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean> | boolean {

    if (this.authService.isTokenExpired()) {
      this.authService.logout()
      this.router.navigateByUrl('/')
      this.alertifyService.error("Oturum Süreniz Sona Erdi!")
      return false
    } else {
      return new Observable<boolean>(obs => {
        this.authService.verify().subscribe(response => {
          if (response.success) {
            this.authService.login(response.data.user)
            obs.next(true)
          } else {
            this.authService.logout()
            this.router.navigateByUrl('/')
            this.alertifyService.error("Geçersiz Token!")
            obs.next(false)
          }
        })
      })
    }
  }

}