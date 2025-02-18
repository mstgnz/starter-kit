import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { Dashboard } from './pages/dashboard/dashboard';
import { Documentation } from './pages/documentation/documentation';
import { Landing } from './pages/landing/landing';
import { Notfound } from './pages/notfound/notfound';
import { AppLayout } from './layout/component/app.layout';

const routes: Routes = [
  {
      path: '',
      component: AppLayout,
      children: [
          { path: '', component: Dashboard },
          { path: 'uikit', loadChildren: () => import('./pages/uikit/uikit.routes') },
          { path: 'documentation', component: Documentation },
          { path: 'pages', loadChildren: () => import('./pages/pages.routes') }
      ]
  },
  { path: 'landing', component: Landing },
  { path: 'notfound', component: Notfound },
  { path: 'auth', loadChildren: () => import('./pages/auth/auth.routes') },
  { path: '**', redirectTo: '/notfound' }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class ThemeRoutingModule { }