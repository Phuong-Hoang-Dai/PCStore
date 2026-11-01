import { useSelector } from "react-redux";
import { selectUser } from "../auth/userReducer";
import { Navigate } from "react-router-dom";
import FormSignIn from "../auth/FormSignIn";
import { isAdminRole } from "../auth/user_model";

const Login = () => {
  const user = useSelector(selectUser);

  return isAdminRole(user.roleId) ? (
    <Navigate to="/manage" />
  ) : (
    <>
      <div className="w-full h-screen bg-gray-100 flex justify-center items-center">
        <FormSignIn IsForUser={false} />
      </div>
    </>
  );
};

export default Login;
