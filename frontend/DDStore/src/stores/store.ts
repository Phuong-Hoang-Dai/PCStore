import { configureStore } from "@reduxjs/toolkit";
import cartReducer from "../features/order/cartSlice";
import userReducer from "../features/auth/userReducer";

export const store = configureStore({
  reducer: {
    cart: cartReducer,
    user: userReducer
  },
});


export type AppState = typeof store;
export type RootState = ReturnType<AppState["getState"]>;
export default store;