import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';

import { Router } from '@angular/router';
import { User } from '../graphql/users.graphql';
import { environment } from "../../environments/environment";
import { Permission } from '../interfaces/permission.interface';

@Injectable()
export class AuthService {

    private baseUrl: string = environment.apiEndpoint;
    public currentUser: User = {} as User

    constructor(
        private router: Router,
        private http: HttpClient
    ) { }

    apiLogin(email_or_phone: String, password: String, send: string): Observable<any> {
        return this.http.post<any>(this.baseUrl + "/login", { email_or_phone: email_or_phone, password: password, verify_method: send })
    }

    // token verify with laravel api
    verify(): Observable<any> {
        const options = {
            headers: new HttpHeaders({
                'Authorization': 'Bearer ' + String(this.getToken())
            })
        }
        return this.http.get<any>(this.baseUrl + '/verify', options)
    }

    // logout
    logout(): void {
        this.currentUser = {} as User
        localStorage.removeItem("access_token")
        this.router.navigateByUrl('/')
    }

    // login
    login(user: User): void {
        this.currentUser = user
    }

    getToken() {
        return localStorage.getItem('access_token')
    }

    parseToken() {
        return JSON.parse(atob(this.getToken()!.split('.')[1]))
    }

    isTokenExpired() {
        const token = this.parseToken()
        if (token && token.exp) {
            return Math.floor((new Date).getTime() / 1000) >= token.exp
        }
        return true
    }

    setUserAndModule(): Observable<any> {
        return this.http.post(environment.hasuraHttpEndpoint, {
            query: `query AUTH_USER($id: Int! = ${this.currentUser.id}) {
                users_by_pk(id: $id) {
                    id
                    user_type_id
                    address_id
                    company_id
                    permission_profile_id
                    photo_id
                    firstname
                    lastname
                    email
                    identity_no
                    phone
                    code
                    created_at
                    updated_at
                    deleted_at
                    last_login
                  }
              }`
        }, { headers: { authorization: `Bearer ${this.getToken()}` } })
    }

    permission(link: string, componentName: String): Permission {
        const perms: Permission = {} as Permission
        perms.write = false
        perms.read = false
        perms.active = false
        /* const userModule = this.userModules.find(um => um.module.link == link)
        if (userModule && userModule.permission_profile.permission_profile_menus.length) {
            const find = userModule.permission_profile.permission_profile_menus.find(pm => pm.menu.module_id == userModule.module_id && pm.menu.component == componentName)
            if (find) {
                perms.read = find.read
                perms.write = find.write
                perms.active = find.active
            }
        } */
        return perms
    }

}