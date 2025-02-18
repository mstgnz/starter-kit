import { Observable } from "rxjs";
import { Injectable } from "@angular/core";
import { FormGroup } from "@angular/forms";
import { DateHelper } from "../helpers/date.helper";
import { HttpClient, HttpHeaders } from "@angular/common/http";
import { environment } from "../../environments/environment.development";

@Injectable()
export class ApiService {

    public options = {
        headers: new HttpHeaders({
            'access-token': ""
        })
    }

    constructor(
        private http: HttpClient
    ) { }

    setToken() {
        const token = localStorage.getItem('access_token')
        if (token) {
            if (this.options.headers.has('access-token')) {
                this.options.headers = this.options.headers.set('access-token', token)
            } else {
                this.options.headers = this.options.headers.append('access-token', token)
            }
        }
    }

    profile(): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + "user", this.options)
    }

    verifyCode(email_or_phone: string, code: number): Observable<any> {
        return this.http.post<any>(environment.apiEndpoint + "/verify-code", { email_or_phone: email_or_phone, code: code })
    }

    forgotPassword(email: string, send: string): Observable<any> {
        const formData = new FormData()
        formData.append('email', email)
        formData.append('send', send)
        return this.http.post<any>(environment.apiEndpoint + "user/forgot-password", formData, this.options)
    }

    forgotPasswordChange(email: string, forgot_code: string, password: string, rePassword: string): Observable<any> {
        const formData = new FormData()
        formData.append('email', email)
        formData.append('forgot_code', forgot_code)
        formData.append('password', password)
        formData.append('re_password', rePassword)
        return this.http.post<any>(environment.apiEndpoint + "user/change-forgot-password", formData, this.options)
    }

    changePassword(password: string, rePassword: string): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('password', password)
        formData.append('re_password', rePassword)
        return this.http.post<any>(environment.apiEndpoint + "user/change-password", formData, this.options)
    }

    adminChangePassword(id: Number, password: string, rePassword: string): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('id', String(id))
        formData.append('password', password)
        formData.append('re_password', rePassword)
        return this.http.post<any>(environment.apiEndpoint + "user/admin-change-password", formData, this.options)
    }

    register(firstname: string, lastname: string, email: string, password: string, company_id: string): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('firstname', firstname)
        formData.append('lastname', lastname)
        formData.append('email', email)
        formData.append('password', password)
        formData.append('company_id', company_id)
        return this.http.post<any>(environment.apiEndpoint + "user/register", formData, this.options)
    }

    verify(): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + "user/verify", this.options)
    }

    anonymousToken(): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + "user/anonymous", this.options)
    }

    event(assistant_service_id: Number): Observable<any> {
        this.setToken()
        // AssistantMailEvent, TractorSmsEvent, CheckServiceEvent, ReservationEvent, KmCalculateEvent, limitCheck
        return this.http.get<any>(environment.apiEndpoint + `tursys/event?assistant_service_id=${assistant_service_id}`, this.options)
    }

    manuelProviderAssign(assistant_service_id: Number, provider_reservation_reason_id: Number): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('assistant_service_id', String(assistant_service_id))
        formData.append('provider_reservation_reason_id', String(provider_reservation_reason_id))
        return this.http.post<any>(environment.apiEndpoint + `assistants/manuel-provider-assign`, formData, this.options)
    }

    pricingCalculate(assistant_id: Number): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + `pricing/calculate?assistant_id=${assistant_id}`, this.options)
    }

    nearestLocations(latitude: Number, longitude: Number): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + `companies/nearest-locations?latitude=${latitude}&longitude=${longitude}`, this.options)
    }

    companyGroups(companyId: Number, group_id: Number | string = ""): Observable<any> {
        this.setToken()
        group_id = group_id ? "?=" + group_id : ""
        return this.http.get<any>(environment.apiEndpoint + `companies/groups/${companyId}${group_id}`, this.options)
    }

    reservationCalculate(assistant_service_id: Number): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + `assistant_service/reservation/calculate?assistant_service_id=${assistant_service_id}`, this.options)
    }

    limitCheck(assistant_service_id: Number, update: boolean): Observable<any> {
        this.setToken()
        const link = update ? `assistants/limit-check?assistant_service_id=${assistant_service_id}&update=true` : `assistants/limit-check?assistant_service_id=${assistant_service_id}`
        return this.http.get<any>(environment.apiEndpoint + link, this.options)
    }

    incidentCreate(address: string, district: string, city: string, latitude: number, longitude: number): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('address', address)
        formData.append('city', city)
        formData.append('district', district)
        formData.append('latitude', String(latitude))
        formData.append('longitude', String(longitude))
        return this.http.post<any>(environment.apiEndpoint + "address/incident-create", formData, this.options)
    }

    smsAddress(phone: string, assistant_service_id: Number): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('phone', phone)
        formData.append('assistant_service_id', String(assistant_service_id))
        return this.http.post<any>(environment.apiEndpoint + "sms/incident", formData, this.options)
    }

    assistantUpdate(formGroup: FormGroup): Observable<any> {
        this.setToken()
        const formData = new FormData()
        for (const [key, value] of Object.entries(formGroup.value)) {
            if (typeof value != "object" && value != null) {
                formData.append(key, String(value))
            }
        }
        return this.http.post<any>(environment.apiEndpoint + "assistants/update", formData, this.options)
    }

    assistantServiceCreate(formGroup: FormGroup): Observable<any> {
        this.setToken()
        const formData = new FormData()
        for (const [key, value] of Object.entries(formGroup.value)) {
            if (typeof value != "object" && value != null) {
                formData.append(key, String(value))
            }
        }
        return this.http.post<any>(environment.apiEndpoint + "assistants/services/create", formData, this.options)
    }

    assistantServiceUpdate(formGroup: FormGroup): Observable<any> {
        this.setToken()
        const formData = new FormData()
        for (const [key, value] of Object.entries(formGroup.value)) {
            if (typeof value != "object" && value != null) {
                formData.append(key, String(value))
            }
        }
        return this.http.post<any>(environment.apiEndpoint + "assistants/services/update", formData, this.options)
    }

    policyControl(firma: string, kayitno: string): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('username', "kferhat")
        formData.append('password', "x%G@bNtYm#{8")
        formData.append('firma', firma)
        formData.append('kayitno', kayitno)
        return this.http.post<any>("https://app.turassist.com/policekontrol/jsonapiservice.php", formData, this.options)
    }

    sendMail(to: string, subject: string, content: string, bcc: string[] = []): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('to', to)
        formData.append('subject', subject)
        formData.append('content', content)
        if (bcc.length) {
            bcc.forEach(b => {
                formData.append('bcc[]', b)
            })
        }
        return this.http.post<any>(environment.apiEndpoint + "mail/send", formData, this.options)
    }

    rentGoGet(download: boolean = false): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + (download ? 'rent-go/get?download' : 'rent-go/get'), this.options)
    }

    rentGoPull(): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + 'rent-go/pull', this.options)
    }

    rentGoFind(plateNumber: string): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + `rent-go/find?plate_number=${plateNumber}`, this.options)
    }

    nissanFind(plateNumber: string = "", chassisNo: string = ""): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + `bring/nissan_portfolio?plate_number=${plateNumber}&chassis_no=${chassisNo}`, this.options)
    }

    sendPayLink(assistantServiceId: Number, phone: Number): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('phone', String(phone))
        formData.append('assistant_service_id', String(assistantServiceId))
        return this.http.post<any>(environment.apiEndpoint + "payment/token", formData, this.options)
    }

    costList(page: number, all: boolean): Observable<any> {
        this.setToken()
        let link = all ? `assistant_service/cost-list?page=${page}&all` : `assistant_service/cost-list?page=${page}`
        return this.http.get<any>(environment.apiEndpoint + link, this.options)
    }

    nearestPlace(service_id: Number, latitude: Number, longitude: Number): Observable<any> {
        this.setToken()
        let params = `?service_id=${service_id}&latitude=${latitude}&longitude=${longitude}`
        return this.http.get<any>(environment.apiEndpoint + 'assistant_service/nearest-place' + params, this.options)
    }

    audatex(assistant_service_id: Number): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('assistant_service_id', String(assistant_service_id))
        return this.http.post<any>(environment.apiEndpoint + "audatex/image-capture", formData, this.options)
    }

    audatexCreateTask(assistant_service_id: Number): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + `audatex/create-task?assistant_service_id=${assistant_service_id}`, this.options)
    }

    audatexGetTask(task_id: string): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + `audatex/get-task?task_id=${task_id}`, this.options)
    }

    audatexSendExper(insurer_code: string, assessor_code: string, file_number: string, xml: string): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('insurer_code', insurer_code)
        formData.append('assessor_code', assessor_code)
        formData.append('file_number', file_number)
        formData.append('xml', xml)
        return this.http.post<any>(environment.apiEndpoint + "audatex/send-exper", formData, this.options)
    }

    useRegionAllCountry(company_id: Number, company_assistant_region_id: Number): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('company_id', String(company_id))
        formData.append('company_assistant_region_id', String(company_assistant_region_id))
        return this.http.post<any>(environment.apiEndpoint + `companies/use-region-all-region`, formData, this.options)
    }

    geocode(address: string): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('address', address)
        return this.http.post<any>(environment.apiEndpoint + "map/geocode", formData, this.options)
    }

    kmCalculate(assistant_service_id: Number): Observable<any> {
        this.setToken()
        let params = `?assistant_service_id=${assistant_service_id}`
        return this.http.get<any>(environment.apiEndpoint + "map/km-calculate" + params, this.options)
    }

    nissanXml(start_date: string, end_date: string): Observable<any> {
        this.setToken()
        return this.http.get<any>(environment.apiEndpoint + `webservices/nissan-xml?start_date=${start_date}&end_date=${end_date}`, this.options)
    }

    toyotaList(startDate: Date, endDate: Date, isActive: boolean): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('startDate', DateHelper.dateFormat(String(startDate)))
        formData.append('endDate', DateHelper.dateFormat(String(endDate)))
        formData.append('isActive', String(isActive))
        return this.http.post<any>(environment.apiEndpoint + "webservices/toyota-list", formData, this.options)
    }

    toyotaGet(fullName: string = "", plateNumber: string = "", vinNumber: string = "", vinNumberLast7Char: string = ""): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('fullName', fullName)
        formData.append('plateNumber', plateNumber)
        formData.append('vinNumber', vinNumber)
        formData.append('vinNumberLast7Char', vinNumberLast7Char)
        return this.http.post<any>(environment.apiEndpoint + "webservices/toyota-get", formData, this.options)
    }

    uttsForm(formGroup: FormGroup, path: string = "create"): Observable<any> {
        this.setToken()
        const formData = new FormData()
        for (const [key, value] of Object.entries(formGroup.value)) {
            if (typeof value != "object" && value != null) {
                formData.append(key, String(value))
            }
        }
        return this.http.post<any>(environment.apiEndpoint + "utts/" + path, formData, this.options)
    }

    uttsExcel(formGroup: FormGroup, file: File): Observable<any> {
        this.setToken()
        const formData = new FormData()
        for (const [key, value] of Object.entries(formGroup.value)) {
            if (typeof value != "object" && value != null) {
                formData.append(key, String(value))
            }
        }
        formData.append("file", file)
        return this.http.post<any>(environment.apiEndpoint + "utts/excel", formData, this.options)
    }

    uttsCode(firstname: string, lastname: string, phone: string): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('firstname', firstname)
        formData.append('lastname', lastname)
        formData.append('phone', phone)
        return this.http.post<any>(environment.apiEndpoint + "utts/code", formData, this.options)
    }

    smsOpen(link: string): Observable<any> {
        this.setToken()
        const formData = new FormData()
        formData.append('link', link)
        return this.http.post<any>(environment.apiEndpoint + "sms/open", formData, this.options)
    }

    fetchXmlContent(xmlUrl: string): Observable<string> {
        return this.http.get(xmlUrl, { responseType: 'text' });
    }

}