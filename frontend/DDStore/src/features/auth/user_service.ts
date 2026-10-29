import type { User, UserLogin } from "./user_model";

const url = import.meta.env.VITE_API_URL + "/user"

const emptyUser: User = {
  id: 0,
  name: "",
  email: "",
  roleId: "",
}

export const SignUp = async (user: UserLogin) => {
  try{
    const res = await fetch(url, {
      method: "POST",
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(user)
    })
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
    const res = await fetch(url + "/login", {
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
    const res = await fetch(url + "/logout", {
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
    const res = await fetch(url + "/me", {
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
