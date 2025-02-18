import { HttpClient } from "@angular/common/http";
import { Injectable } from "@angular/core";

@Injectable()
export class SoundService {

    constructor(
        private http: HttpClient
    ) { }

    newOrder() {
        new Audio("assets/sounds/notification.wav").play()
    }


}