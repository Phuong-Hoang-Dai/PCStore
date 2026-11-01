
export type UserLogin = {
  id: number
  email:string
  name: string
  password: string
}

export type User = {
  id:number
  email: string
  name:string
  roleId: string
}

// Roles the gateway lets into admin routes.
export const isAdminRole = (role: string) => role === "admin" || role === "owner"