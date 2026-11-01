import { CiShoppingCart } from "react-icons/ci";
import { selectCartTotal } from "./cartSlice";
import { useSelector } from "react-redux";
import { useState, useRef, useEffect } from "react";
import CartList from "./CartList";

const Cart = () => {
  const quantity: number = useSelector(selectCartTotal);
  const refCart = useRef<HTMLDivElement>(null);
  const [isShow, setIsShow] = useState(false);
  const handleClick = (e: React.MouseEvent<HTMLButtonElement, MouseEvent>) => {
    e.preventDefault();
    setIsShow((prev) => !prev);
  };
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      if (!refCart.current?.contains(e.target as Node)) {
        setIsShow(false);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, []);
  return (
    <div className="relative flex" ref={refCart}>
      <button
        className="relative flex text-6xl cursor-pointer"
        onClick={(e) => handleClick(e)}
      >
        <CiShoppingCart size={45} />
        <span className="text-white text-xs bg-red-500 px-2 py-1 rounded-4xl absolute right-0 top-0">
          {quantity}
        </span>
      </button>
      {isShow && (
        <div className="absolute top-20 right-0 z-10 cursor-default">
          <CartList onClose={() => setIsShow(false)} />
        </div>
      )}
    </div>
  );
};

export default Cart;
