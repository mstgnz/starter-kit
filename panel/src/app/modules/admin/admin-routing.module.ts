import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NotfoundComponent } from '../lobby/notfound.component';

const routes: Routes = [
  {
    path: '', children: [
      {
        path: '', data: { breadcrumb: 'Anasayfa', name: 'DashboardComponent' },
        component: DashboardComponent, pathMatch: 'full'
      },
    ]
  },
  { path: 'notfound', component: NotfoundComponent },
  { path: '**', redirectTo: '/notfound' }
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule]
})
export class AdminRoutingModule { }