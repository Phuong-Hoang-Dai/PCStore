import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import {type RootState  } from "../../stores/store";
import type { Product } from "../product/product_model";

interface Cart {
  order: Product[]
  total: number
}

const initialState: Cart = {
  order: [],
  total: 0,
};

export const cartCount = createSlice({
  name: "cart",
  initialState,
  reducers: {
    addItem:(state: Cart, action: PayloadAction<Product>) => {
      state.total += action.payload.quantityOrder
      const existingItem = state.order.find((product) => product.id === action.payload.id)
      if (existingItem){
        existingItem.quantityOrder += action.payload.quantityOrder
      }else{
        state.order.push(action.payload)
      }
    },
    increment: (state: Cart, action: PayloadAction<Product>) => {
      state.total += 1;
      const existingItem = state.order.find((product) => product.id === action.payload.id)
      if (existingItem){
        existingItem.quantityOrder += 1
      }else{
        console.log({ ...action.payload, quantityOrder: 1 });
        state.order.push({ ...action.payload, quantityOrder: 1 });
      }
    },
    decrement: (state: Cart, action: PayloadAction<Product>) => {
      state.total -= 1
      if (state.total < 0){
        state.total = 0
      }
      const existingItem = state.order.find((product) => product.id === action.payload.id)
      if (existingItem) {
        existingItem.quantityOrder -= 1;
      }
      if (existingItem && existingItem.quantityOrder <= 0){
          state.order = state.order.filter(product => product.id !== existingItem.id)
      }
    },
    changeByAmout: (state: Cart, action: PayloadAction<{item: Product, quantity: number}>) => {
      state.total += action.payload.quantity;
      if (state.total < 0) {
        state.total = 0;
      }
          const existingItem = state.order.find((product) => product.id === action.payload.item.id)
      if (existingItem) {
        existingItem.quantityOrder -= action.payload.quantity;
      }
      if (existingItem && existingItem.quantityOrder <= 0){
          state.order = state.order.filter(product => product.id !== existingItem.id)
      }
    },
  },
});

export const {addItem, increment, decrement, changeByAmout } = cartCount.actions;
export default cartCount.reducer;
export const selectCartTotal = (state: RootState) => state.cart.total;
export const selectCartItems  =(id:number) => (state: RootState) => state.cart.order.find(product => product.id === id);
export const selectCart= (state: RootState) => state.cart;
