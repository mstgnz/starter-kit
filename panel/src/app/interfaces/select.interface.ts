export interface SelectData {
    id: Number
    name: String
}

export interface StepData {
    step: number
    next: boolean
    update?: boolean
}

export interface UnfinishedSearch {
    company_id: number
    contract_id: number
    call_reason_id: number
    assistant_type_id: number
    representative_firstname: string
    representative_lastname: string
    customer_firstname: string
    customer_lastname: string
    customer_phone: string
    plate_number: string
    policy_number: string
}

export interface DateTimeColumn {
    date: string[]
    time: string[]
}