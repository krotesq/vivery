export interface AccountMe {
  id: string
  username: string
  createdAt: string
}

export interface Account {
  id: string
  username: string
  active: boolean
  failedLoginAttempts: number
  lockedUntil: string
  createdAt: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface CreateAccountRequest {
  username: string
  password: string
}