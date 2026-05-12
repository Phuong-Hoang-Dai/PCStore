import { createSlice, type PayloadAction } from "@reduxjs/toolkit"
import type { User } from "./user_model"
import type { RootState } from "../../stores/store";

interface UserInfo {
  user: User
}

const initialState: UserInfo = {
  user: {
    id: 0,
    name: "",
    email: "",
    roleId: ""
  },
};

export const userReducer = createSlice({
  name: "user",
  initialState,
  reducers: {
    SetUser: (state: UserInfo, payload: PayloadAction<User>) => {
      state.user.id = payload.payload.id
      state.user.name = payload.payload.name
      state.user.email = payload.payload.email
      state.user.roleId = payload.payload.roleId
    },
    RemoveUser: (state: UserInfo) => {
      state.user.id = 0
      state.user.name = ""
      state.user.email = ""
      state.user.roleId = ""
    }
  }
})

export const { SetUser,RemoveUser } = userReducer.actions
export default userReducer.reducer
export const selectUser = (state: RootState) => state.user.user