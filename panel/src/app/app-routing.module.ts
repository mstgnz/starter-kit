import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';

import { AppGuard } from './guards/app.guard';
import { AppLayout } from './layout/component/app.layout';
import { LoginComponent } from './modules/lobby/login.component';
import { NotfoundComponent } from './modules/lobby/notfound.component';
import { NewPasswordComponent } from './modules/lobby/new-password.component';
import { ForgotPasswordComponent } from './modules/lobby/forgot-password.component';

const routes: Routes = [
  {
      path: '',
      component: AppLayout,
      children: [
          { path: 'admin', loadChildren: () => import('./modules/panel/panel.module').then(m => m.PanelModule), canActivate: [AppGuard] },
          { path: 'forgot-password', component: ForgotPasswordComponent },
          { path: 'new-password', component: NewPasswordComponent },
      ]
  },
  { path: 'theme', loadChildren: () => import('./modules/theme/theme.module').then(m => m.ThemeModule) },
  { path: 'notfound', component: NotfoundComponent },
  { path: 'login', component: LoginComponent },
  { path: '**', redirectTo: '/notfound' }
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
