import { useEffect, useState } from "react";
import { selectUser, SetUser } from "./userReducer";
import { useDispatch, useSelector } from "react-redux";
import { Navigate } from "react-router-dom";
import { Outlet } from "react-router-dom";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import { isAdminRole } from "./user_model";
import { TryGetUserByCookie } from "./user_service";

const RequreAdmin = () => {
  const user = useSelector(selectUser);
  const dispatch = useDispatch();
  // The admin pages have no NavBar to restore the session, so on a page refresh the
  // cookie is checked here before deciding to redirect.
  const [isChecking, setIsChecking] = useState(user.id === 0);

  useEffect(() => {
    if (!isChecking) {
      return;
    }
    TryGetUserByCookie()
      .then((data) => dispatch(SetUser(data)))
      .finally(() => setIsChecking(false));
  }, [isChecking, dispatch]);

  if (isChecking) {
    return (
      <div className="w-full h-screen flex justify-center items-center">
        <AiOutlineLoading3Quarters className="animate-spin" size={32} />
      </div>
    );
  }

  return isAdminRole(user.roleId) ? <Outlet /> : <Navigate to="/admin/login" />;
};

export default RequreAdmin;
