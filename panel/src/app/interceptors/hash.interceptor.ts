/// <reference lib="dom" />

import { from } from 'rxjs';
import { inject } from '@angular/core';
import { switchMap } from 'rxjs/operators';
import { HashService } from '../services/hash.service';
import { environment } from '../../environments/environment';
import { HttpInterceptorFn, HttpRequest, HttpHandlerFn } from '@angular/common/http';

const getPathFromUrl = (url: string): string => {
    const baseUrl = environment.apiEndpoint;
    return url.replace(baseUrl, '');
};

export const hashInterceptor: HttpInterceptorFn = (req: HttpRequest<unknown>, next: HttpHandlerFn) => {
    const hashService = inject(HashService);

    // Sadece API isteklerini yakala
    if (!req.url.startsWith(environment.apiEndpoint)) {
        return next(req);
    }

    // Path'i al
    const path = getPathFromUrl(req.url);

    // Hash oluştur ve header'ları ekle
    return from(hashService.setHeaderHash(path)).pipe(
        switchMap(hash => {
            // Mevcut header'ları koru ve yenilerini ekle
            let headers = req.headers
                .set('Hash', hash.Hash)
                .set('Timestamp', hash.Timestamp);

            // Token varsa ekle
            const token = localStorage.getItem('access_token');
            if (token) {
                headers = headers.set('Authorization', `Bearer ${token}`);
            }

            // Request'i klonla ve header'ları ekle
            const modifiedRequest = req.clone({
                headers: headers
            });

            return next(modifiedRequest);
        })
    );
};