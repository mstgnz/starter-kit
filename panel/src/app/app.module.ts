import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import Aura from '@primeng/themes/aura';
import { AppGuard } from './guards/app.guard';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AppComponent } from './app.component';
import { providePrimeNG } from 'primeng/config';
import { GraphQLModule } from './graphql.module';
import { LoginGuard } from './guards/login.guard';
import { ApiService } from './services/api.service';
import { CdnService } from './services/cdn.service';
import { SmsService } from './services/sms.service';
import { SharedModule } from './shared/shared.module';
import { AuthService } from './services/auth.service';
import { AppRoutingModule } from './app-routing.module';
import { ExcelService } from './services/excel.service';
import { TimerService } from './services/timer.service';
import { SoundService } from './services/sound.service';
import { GoogleService } from './services/google.service';
import { PermissionGuard } from './guards/permission.guard';
import { DownloadService } from './services/download.service';
import { AlertifyService } from './services/alertify.service';
import { LoginComponent } from './modules/lobby/login.component';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { hashInterceptor } from './interceptors/hash.interceptor';
import { NotfoundComponent } from './modules/lobby/notfound.component';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { httpErrorInterceptor } from './interceptors/http-error.interceptor';
import { NewPasswordComponent } from './modules/lobby/new-password.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { AccessDeniedComponent } from './modules/lobby/access-denied.component';
import { ForgotPasswordComponent } from './modules/lobby/forgot-password.component';

@NgModule({
  declarations: [
    AppComponent,
    LoginComponent,
    NotfoundComponent,
    NewPasswordComponent,
    AccessDeniedComponent,
    ForgotPasswordComponent
  ],
  imports: [
    FormsModule,
    CommonModule,
    RouterModule,
    SharedModule,
    BrowserModule,
    GraphQLModule,
    AppRoutingModule,
    ReactiveFormsModule,
    BrowserAnimationsModule,
  ],
  providers: [
    AppGuard,
    LoginGuard,
    ApiService,
    CdnService,
    SmsService,
    AuthService,
    ExcelService,
    TimerService,
    SoundService,
    GoogleService,
    DownloadService,
    PermissionGuard,
    AlertifyService,
    provideHttpClient(
      withInterceptors([hashInterceptor, httpErrorInterceptor])
    ),
    providePrimeNG({ theme: { preset: Aura, options: { darkModeSelector: '.app-dark' } } })
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
