/// <reference lib="dom" />

import { throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { HttpInterceptorFn } from '@angular/common/http';

export const httpErrorInterceptor: HttpInterceptorFn = (req, next) => {
    return next(req).pipe(
        catchError(error => {
            let message = "Hata Oluştu!";

            if (!navigator.onLine) {
                message = "İnternet bağlantınız yok";
                return throwError(() => new Error(message));
            }

            if (error.error?.error) {
                if (error.status === 401) {
                    message = "Yetkiniz Yok!";
                    return throwError(() => new Error(message));
                }

                switch (error.error.error.message) {
                    case "EMAIL_EXISTS":
                        message = "Email Adresi Zaten Kayıtlı!";
                        break;
                    case "EMAIL_NOT_FOUND":
                        message = "Email Adresi Bulunamadı!";
                        break;
                    case "INVALID_PASSWORD":
                        message = "Parolanız Yanlış!";
                        break;
                    case "USER_DISABLED":
                        message = "Kullanıcı Aktif Değil!";
                        break;
                    case "OPERATION_NOT_ALLOWED":
                        message = "Kullanıcı Girişi Kapalı!";
                        break;
                    case "TOO_MANY_ATTEMPTS_TRY_LATER":
                        message = "Bu Cihaz ile Girişler Kapatıldı!";
                        break;
                    default:
                        message = error.error.error.message;
                        break;
                }
            }
            return throwError(() => new Error(message));
        })
    );
};