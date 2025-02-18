import { Injectable } from "@angular/core";
import { gql, Mutation, Query } from "apollo-angular";


export interface User {
  id: Number
  user_type_id: Number
  address_id: Number
  company_id: Number
  permission_profile_id: Number
  photo_id: Number
  firstname: String
  lastname: String
  email: String
  identity_no: Number
  phone: String
  code: String
  created_at: Date
  updated_at: Date
  deleted_at: Date
  last_login: Date
}

interface ResponseArray {
  users: User[]
}

interface ResponseSingle {
  users_by_pk: User
}

interface Aggregate {
  users_aggregate: {
    aggregate: {
      count: number
    }
  }
}

interface InsertMutation {
  insert_users: {
    affected_rows: number
    returning: [{
      id: Number
    }]
  }
}

interface UpdateMutation {
  update_users: {
    affected_rows: number
    returning: [{
      id: Number
    }]
  }
}

interface DeleteMutation {
  delete_users: {
    affected_rows: number
    returning: [{
      id: Number
    }]
  }
}

@Injectable()
export class USER_COUNT extends Query<Aggregate> {
  override document = gql`
  query USER_COUNT($where: users_bool_exp) {
    users_aggregate(where: $where) {
      aggregate {
        count
      }
    }
  }`;
}

@Injectable()
export class USERS extends Query<ResponseArray> {
  override document = gql`
  query USERS($where: users_bool_exp){
    users(where:$where){
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
  }`;
}

@Injectable()
export class USER_CREATE extends Mutation<InsertMutation> {
  override document = gql`
  mutation USER_CREATE($company_id:Int!, $email:String! $active:Boolean!, $priority:smallint!) {
    insert_users(objects: {company_id:$company_id, email:$email, active:$active, priority:$priority}) {
      affected_rows
      returning {
        id
      }
    }
  }`;
}

@Injectable()
export class USER_UPDATE extends Mutation<UpdateMutation> {
  override document = gql`
  mutation USER_UPDATE($id: Int!, $company_id:Int!, $email:String!, $active:Boolean!, $priority:smallint!) {
    update_users(where: {id: {_eq: $id}}, _set: {company_id:$company_id, email:$email, active:$active, priority:$priority}) {
      affected_rows
      returning {
        id
      }
    }
  }`;
}

@Injectable()
export class USER_DELETE extends Mutation<DeleteMutation> {
  override document = gql`
  mutation USER_DELETE($id: Int!) {
    delete_users(where: {id: {_eq: $id}}) {
      affected_rows
      returning {
        id
      }
    }
  }`;
}