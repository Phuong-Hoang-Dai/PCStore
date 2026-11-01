import { configureStore } from "@reduxjs/toolkit";
import cartReducer, { CART_STORAGE_KEY } from "../features/order/cartSlice";
import userReducer from "../features/auth/userReducer";

export const store = configureStore({
  reducer: {
    cart: cartReducer,
    user: userReducer
  },
});

let savedCart = store.getState().cart;
store.subscribe(() => {
  const cart = store.getState().cart;
  if (cart === savedCart) {
    return;
  }
  savedCart = cart;
  try {
    localStorage.setItem(CART_STORAGE_KEY, JSON.stringify(cart));
  } catch {
    // Storage can be unavailable (private mode, quota); the cart still works in memory.
  }
});


export type AppState = typeof store;
export type RootState = ReturnType<AppState["getState"]>;
export default store;
