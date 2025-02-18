import { NgModule } from '@angular/core';

import { PanelRoutingModule } from './panel-routing.module';
import { DashboardComponent } from './dashboard/dashboard.component';

@NgModule({
  declarations: [
    DashboardComponent,
  ],
  imports: [
    PanelRoutingModule,
  ],
  exports: []
})
export class PanelModule { }