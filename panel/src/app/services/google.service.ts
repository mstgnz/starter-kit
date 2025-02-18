import { HttpClient, HttpHeaders } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";

import { LocationPoint } from "../interfaces/map.interface";
import { environment } from "../../environments/environment";

@Injectable()
export class GoogleService {

    private url: string | undefined
    private options = {
        headers: new HttpHeaders({

        })
    }

    constructor(
        private http: HttpClient
    ) { }

    getAddress(lat: Number, lng: Number): Observable<any> {
        this.url = `https://maps.googleapis.com/maps/api/geocode/json?key=${environment.googleMapApiKey}&latlng=${lat},${lng}`
        return this.http.get<any>(this.url, this.options)
    }

    getAddress2(address: string): Observable<any> {
        this.url = `https://maps.googleapis.com/maps/api/geocode/json?key=${environment.googleMapApiKey}&address=${address}`
        return this.http.get<any>(this.url, this.options)
    }

    getDirections(origin: LocationPoint, destination: LocationPoint): Observable<any> {
        this.url = `https://maps.googleapis.com/maps/api/directions/json?key=${environment.googleMapApiKey}&origin=${origin.latitude},${origin.longitude}&destination=${destination.latitude},${destination.longitude}`
        return this.http.get<any>(this.url, this.options)
    }

    getDistance(origin: LocationPoint, destination: LocationPoint): Observable<any> {
        this.url = `https://maps.googleapis.com/maps/api/distancematrix/json?key=${environment.googleMapApiKey}&origins=${origin.latitude},${origin.longitude}&destinations=${destination.latitude},${destination.longitude}`
        return this.http.get<any>(this.url, this.options)
    }
}