import { CiShoppingCart } from "react-icons/ci";
import { selectCartTotal } from "./cartSlice";
import { useSelector } from "react-redux";
import { useState, useRef, useEffect } from "react";
import CartList from "./CartList";

const Cart = () => {
  const quantity: number = useSelector(selectCartTotal);
  console.log("Cart quantity:", quantity);
  const refShowCart = useRef<HTMLDivElement>(null);
  const [isShow, setIsShow] = useState(false);
  const handleClick = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    e.preventDefault();
    setIsShow(true);
  };
  const handleClickOutside = (e: MouseEvent) => {
    if (!refShowCart.current?.contains(e.target as Node)) {
      setIsShow(false);
    }
  };
  useEffect(() => {
    document.addEventListener("mousedown", (e) => handleClickOutside(e));

    return () => {
      document.removeEventListener("mousedown", (e) => handleClickOutside(e));
    };
  }, []);
  return (
    <>
      <button
        className="relative flex text-6xl cursor-pointer"
        onClick={(e) => handleClick(e)}
      >
        <CiShoppingCart size={45} />
        <span className="text-white text-xs bg-red-500 px-2 py-1 rounded-4xl absolute right-0 top-0">
          {quantity}
        </span>
        {isShow && (
          <div
            className="absolute top-20 right-0 z-10 cursor-default"
            ref={refShowCart}
          >
            <CartList />
          </div>
        )}
      </button>
    </>
  );
};

export default Cart;
