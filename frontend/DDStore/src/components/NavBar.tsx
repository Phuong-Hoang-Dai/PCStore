import { Link } from "react-router-dom";
import UserInfo from "../features/user/UserInfo";
import Cart from "../features/order/Cart";

const NavBar = () => {
  return (
    <>
      <div className="bg-[#ecb996] h-15 md:h-20 w-full fixed z-1000">
        <div className="flex flex-row justify-between items-center h-full w-full md:w-9/10 lg:w-17/25 m-auto px-2">
          <Link to="/" className="h-full">
            <img src="/logo.png" alt="DDStore" className="h-full" />
          </Link>
          
          <div className="flex flex-row items-center h-full">
            <UserInfo />
            <Cart />
          </div>
        </div>
      </div>
      <div className="bg-[#ecb996] h-15 md:h-20 w-full"></div>
    </>
  );
};

export default NavBar;
