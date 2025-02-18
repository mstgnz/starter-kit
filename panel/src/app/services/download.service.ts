import { HttpClient, HttpHeaders } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";

@Injectable()
export class DownloadService {

    private options = {
        headers: new HttpHeaders({})
    }

    constructor(
        private http: HttpClient
    ) { }

    download(url: string): Observable<Blob> {
        return this.http.get(url, {
            responseType: 'blob'
        })
    }

}