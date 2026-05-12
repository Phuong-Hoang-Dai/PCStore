import { selectUser } from "./userReducer";
import { useSelector } from "react-redux";
import { Navigate } from "react-router-dom";
import { Outlet } from "react-router-dom";

const RequreAdmin = () => {
  const user = useSelector(selectUser);

  return (
    <>{user.roleId != "admin" ? <Navigate to="/admin/login" /> : { Outlet }}</>
  );
};

export default RequreAdmin;
