export interface CompanyListSearch {
    name: String
    code: Number
    city: String
    district: String
    tax_no: Number
    service_name: String
}

export interface AsistantListSearch {
    id: Number
    tc_no: Number
    policy_number: String
    phone: String
    plate_number: String
    chassis_no: String
    assistant_service_id: Number
    company: Number | String
    provider: Number | String
    service: Number | String
    file_open_date: Date
    file_close_date: Date
    check_policy: Boolean | undefined
}

export interface PolicyListSearch {
    company_name: String
    firstname: String
    lastname: String
    plate_number: String
    policy_number: String
}

export interface CompanyAssistantRegionServiceSearch {
    service_name: String
    company_name: String
    company_code: String
    assistant_type: String
    city: String
    district: String
    priority: Number
    ratio: Number
    group_id: Number
}

export interface ContractListSearch {
    company_name: String
    contract_name: String
    contract_code: String
    contract_type_name: String
}

export interface ContractCallReasonListSearch {
    contract_call_reason_name: String
    company_name: String
    contract_name: String
}

export interface RawDataSearch {
    plaka: String
    markamodel: String
    adsoyad: String
    fileno: String
    il: String
    ilce: String
    hizmettarihi: Date
    start_date: Date
    end_date: Date
    nbcontra: String
    nbcause: String
    tcdescri: String
    musteri_adi: String
    provider_adi: String
    sigortali_adi: String
}

export interface ExgratiaSearch {
    firstname: String
    lastname: String
}

export interface AsistantServiceListSearch {
    id: Number
    tc_no: Number
    phone: String
    chassis_no: String
    plate_number: String
    policy_number: String
    company: Number | String
    provider: Number | String
    service: Number | String
    file_open_date: Date
    file_close_date: Date
    check_policy: Boolean | undefined
}

export interface AlotechListSearch {
    agent: String
    called_num: String
}

export interface ReportSearch {
    company_name: String
    provider_name: String
}

export interface CallCenterQuerySearch {
    firstname: String
    lastname: String
    email: String
}

export interface FollowSearch {
    search: String
    company_name: String
    provider_name: String
    service_name: String
    plate_number: String
    assigner: String
    city_id: Number
    district_id: Number
    assistant_type_id: Number
    assistant_status_id: Number
    assistant_service_id: Number
    assistant_service_status_id: Number
    file_open_date: Date
    file_close_date: Date
    open_appointment_date: Date
    close_appointment_date: Date
}

export interface UserSearch {
    firstname: String
    lastname: String
    email: String
    phone: String
    company_name: String
}