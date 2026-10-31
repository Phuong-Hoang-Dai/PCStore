import type { User, UserLogin } from "./user_model";
import { apiFetch } from "../../utils/api";

const emptyUser: User = {
  id: 0,
  name: "",
  email: "",
  roleId: "",
}

// Thrown by SignUp with a message that can be shown to the user.
export class SignUpError extends Error {}

export const SignUp = async (user: UserLogin) => {
  try{
    const res = await apiFetch("/user", {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(user)
    })
    if (res.status === 409){
      throw new SignUpError("Email này đã có tài khoản, vui lòng đăng nhập");
    }
    if (res.status === 400){
      throw new SignUpError("Họ tên, email hoặc mật khẩu không hợp lệ");
    }
    if (!res.ok){
      throw new Error("Network response was not ok");
    }

    const data = await res.json()
    return data
  }catch(error){
    console.error("Failed to sign up:", error);
    throw error;
  }
}

export const SignIn = async (user: UserLogin) => {
  try{
    const res = await apiFetch("/user/login", {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(user),
      credentials: "include",
    })
    const data = await res.json()
    if (!res.ok){
      throw new Error("Network response was not ok" );
    }
    return data
  }catch(error){
    console.error("Failed to sign in:", error);
    throw error;
  }
}

export const Logout = async () => {
 try{
    const res = await apiFetch("/user/logout", {
      method: "GET",
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: `include`,
    })
    const data = await res.json()
    return data
  }catch(error){
     console.error("Failed to log out:", error);
    throw error;
  }
}

// Returns the empty user (id 0) when there is no valid session cookie.
export const TryGetUserByCookie = async ():Promise<User> => {
  try{
    const res = await apiFetch("/user/me", {
      method: "GET",
      headers: {
        'Content-Type': 'application/json',
      },
      credentials: `include`,
    })
    if (!res.ok){
      return emptyUser
    }
    const data = await res.json()
    const user: User = data["data"]
    return user ?? emptyUser
  }catch(error){
     console.error("Failed to get current user:", error);
    return emptyUser
  }
}
