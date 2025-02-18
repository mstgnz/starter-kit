import { HttpClient, HttpHeaders } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { environment } from "../../environments/environment.development";

@Injectable()
export class SmsService {

    private options = {
        headers: new HttpHeaders({
            'access-token': ""
        })
    }

    constructor(
        private http: HttpClient
    ) { }

    setToken(jwtToken?: string) {
        const token = jwtToken ? jwtToken : localStorage.getItem('access_token')
        if (token) {
            if (this.options.headers.has('access-token')) {
                this.options.headers = this.options.headers.set('access-token', token)
            } else {
                this.options.headers = this.options.headers.append('access-token', token)
            }
        }
    }

    send(header: String, message: String, phones: any[] = [], token?: string): Observable<any> {
        token ? this.setToken(token) : this.setToken()
        const formData = new FormData();
        formData.append('header', String(header));
        formData.append('message', String(message));
        if (phones.length) {
            phones.forEach(phone => {
                formData.append('phones[]', String(phone));
            })
        }
        return this.http.post<any>(environment.apiEndpoint + "sms/send", formData, this.options);
    }

    header(): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + "sms/header");
    }

    credit(): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + "sms/header");
    }

    dlr(dlrId: number): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + "sms/dlr/" + dlrId);
    }

}