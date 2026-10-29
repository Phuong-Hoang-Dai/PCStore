import { FaRegUser } from "react-icons/fa";
import { RemoveUser, selectUser, SetUser } from "../auth/userReducer";
import { useDispatch, useSelector } from "react-redux";
import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import SignUpLogin from "./SignUpLogin";
import { Logout, TryGetUserByCookie } from "../auth/user_service";
import { RiLogoutBoxRLine } from "react-icons/ri";
import { LuReceipt } from "react-icons/lu";

const UserInfo = () => {
  const userInfo = useSelector(selectUser);
  const dispatch = useDispatch();
  const [isPending, setIsPending] = useState(false);
  useEffect(() => {
    TryGetUserByCookie().then((user) => {
      dispatch(SetUser(user));
    });
  }, [dispatch]);

  const handleLogout = (e: React.MouseEvent<SVGElement, MouseEvent>) => {
    e.preventDefault();
    if (!isPending) {
      setIsPending(true);
      Logout()
        .then(() => {
          setIsPending(false);
          dispatch(RemoveUser());
        })
        .catch(() => setIsPending(false));
    }
  };
  return (
    <>
      <div className="h-full">
        {userInfo.id != 0 && (
          <div className="flex flex-row gap-2 justify-end  items-center h-full text-sm font-bold pr-4">
            <Link
              to="/orders"
              className="flex flex-row items-center gap-1 mr-3 hover:underline"
            >
              <LuReceipt />
              <span>Đơn hàng</span>
            </Link>
            <FaRegUser />
            <span className="">{userInfo.name}</span>
            <RiLogoutBoxRLine
              onClick={(e) => handleLogout(e)}
              className="cursor-pointer"
            />
          </div>
        )}
        {userInfo.id == 0 && <SignUpLogin />}
      </div>
    </>
  );
};

export default UserInfo;
