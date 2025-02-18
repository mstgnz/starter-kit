import { HttpClient, HttpHeaders } from "@angular/common/http";
import { Injectable } from "@angular/core";
import { Observable } from "rxjs";
import { environment } from "../../environments/environment.development";

@Injectable()
export class CdnService {

    private options = {
        headers: new HttpHeaders({
            'Authorization': `Bearer ${environment.cdnToken}`
        })
    }

    constructor(
        private http: HttpClient
    ) { }

    upload(pathName: string, file: any, bucketName: string = "turassist"): Observable<any> {
        const formData = new FormData()
        formData.append('file', file)
        formData.append('bucket', bucketName)
        formData.append('path', pathName)
        return this.http.post<any>(environment.cdnEndpoint + "upload", formData, this.options)
    }

    delete(objectName: string, bucketName: string = "turassist"): Observable<any> {
        return this.http.delete<any>(environment.cdnEndpoint + bucketName + "/" + objectName, this.options)
    }

}