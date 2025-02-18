import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, Router, RouterStateSnapshot } from '@angular/router';
import { AuthService } from '../services/auth.service';
import { AlertifyService } from '../services/alertify.service';

@Injectable()
export class PermissionGuard implements CanActivate {

  constructor(
    private router: Router,
    private authService: AuthService,
    private alertifyService: AlertifyService
  ) { }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): boolean {

    if (this.authService.isTokenExpired()) {
      this.authService.logout()
      this.router.navigateByUrl('/')
      this.alertifyService.error("Oturum Süreniz Sona Erdi!")
      return false
    }
    const url = state.url.split('/')[1]
    const componentName = route.data['name']
    if (componentName) {
      const perms = this.authService.permission(url, componentName)
      if (!perms.active) {
        this.alertifyService.error("Sayfa Kullanıma Kapalı! \n" + componentName)
        return false
      }
      if (!perms.read) {
        if (componentName == "DashboardComponent") {
          this.router.navigateByUrl('/')
        } else {
          this.router.navigateByUrl('/' + url)
        }
        this.alertifyService.error("Sayfayı Görüntülüme Yetkiniz Yok! \n" + componentName)
      }
      return true
    }
    this.alertifyService.error("Sayfayı Bulunamadı!")
    return false

  }

}