import React from "react";
import type { Product } from "./product_model";
import { decrement, increment, selectCartItems } from "../order/cartSlice";
import { useDispatch, useSelector } from "react-redux";
import type { UnknownAction } from "@reduxjs/toolkit";
import { BsCartPlus } from "react-icons/bs";
import { LuCircleMinus, LuCirclePlus } from "react-icons/lu";

const UpdateQuantityOrder = ({ product }: { product: Product }) => {
  const quantityOrder = useSelector(selectCartItems(product.id))?.quantityOrder;
  const updateCart = useDispatch();
  const handleButton = (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>,
    action: UnknownAction
  ) => {
    e.preventDefault();
    updateCart(action);
  };

  return (
    <div>
      <div className="h-9 mt-2 flex flex-row">
        <button
          type="button"
          className="cursor-pointer btn pl-2.5 group flex flex-row justify-center items-center relative z-0"
          onClick={(e) => handleButton(e, increment(product))}
        >
          <div
            className={`left-0 w-9 h-full -z-40 absolute  bg-[#FF7601] ${
              quantityOrder ? "" : "group-hover:w-full"
            } transition-all duration-250 rounded-4xl`}
          ></div>
          <BsCartPlus
            className={`text-white text-1xl mr-2  ${
              quantityOrder ? "" : "group-hover:ml-2"
            } transition-all duration-300`}
          />
          {!quantityOrder && (
            <span className="group-hover:text-white transition duration-300 ml-3 md:pr-5 uppercase">
              <span className="hidden md:inline-block">Thêm</span> vào giỏ
            </span>
          )}
        </button>
        {quantityOrder && (
          <div className="ml-4 flex flex-row justify-around items-center w-full text-1xl border border-gray-300 rounded-2xl">
            <button
              className="cursor-pointer hover:bg-gray-100 w-1/4 h-full pr-2 flex items-center justify-end border-r border-gray-300"
              onClick={(e) => handleButton(e, decrement(product))}
            >
              <LuCircleMinus />
            </button>
            <span className="w-2/4 text-center">{quantityOrder}</span>
            <button
              className="cursor-pointer hover:bg-gray-100 w-1/4 h-full flex items-center justify-start pl-2 border-l border-gray-300"
              onClick={(e) => handleButton(e, increment(product))}
            >
              <LuCirclePlus />
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default UpdateQuantityOrder;
