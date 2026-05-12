import { useSelector } from "react-redux";
import { selectUser } from "../auth/userReducer";
import { Navigate } from "react-router-dom";
import FormSignIn from "../auth/FormSignIn";

const Login = () => {
  const user = useSelector(selectUser);

  return user.roleId == "admin" ? (
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
