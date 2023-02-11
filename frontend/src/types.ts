import { Record as PbRecord } from 'pocketbase';

export interface VirtualWallet extends PbRecord {
  balance: number
}

export interface UserDetails extends PbRecord {
  student_id: string
  sex: string
  college_department: string
}

export interface Gift {
  id: string
  uid: string
  label: string
  price: number
}

// NOTE: snake_case because JSON response is in snake_case
interface UserConnection {
  provider: string
}

export interface User extends PbRecord {
  username: string
  email: string
  avatar: string | null
  details: string
  expand: PbRecord['expand'] & {
    wallet: VirtualWallet
    details: UserDetails
  }
  // user_connections: UserConnection[]
}

export interface CollegeDepartment {
  id: string
  label: string
  uid: string
}