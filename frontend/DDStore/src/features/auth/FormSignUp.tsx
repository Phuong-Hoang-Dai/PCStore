import { useState } from "react";
import { AiOutlineLoading3Quarters } from "react-icons/ai";

import { SignIn, SignUp, TryGetUserByCookie } from "./user_service";
import { useDispatch } from "react-redux";
import type { UserLogin } from "./user_model";
import { SetUser } from "./userReducer";

const FormSignUp = () => {
  const dispatch = useDispatch();
  const [form, setForm] = useState({
    id: 0,
    name: "",
    email: "",
    password: "",
    confirmPass: "",
  });
  const [isPending, setIsPending] = useState(false);
  const labels = ["Họ tên", "Email", "Mật khẩu", "Nhập lại mật khẩu"];

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsPending(true);
    const user: UserLogin = {
      id: form.id,
      name: form.name,
      email: form.email,
      password: form.password,
    };

    SignUp(user)
      .then(() => {
        setIsPending(false);
        return SignIn(user);
      })
      .then(() => TryGetUserByCookie())
      .then((data) => dispatch(SetUser(data)))
      .catch();
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    key: string
  ) => {
    e.preventDefault();
    setForm((prev) => ({ ...prev, [key]: e.target.value }));
  };

  const [, ...keys] = Object.keys(form);

  const renderForm = keys.map((key, i) => (
    <div className="mb-4" key={key}>
      <label className="block mb-1 font-medium">{labels[i]}</label>
      <input
        type="text"
        name={key}
        required
        onChange={(e) => handleChange(e, key)}
        className="w-full border border-gray-300 p-2 rounded-lg"
      />
    </div>
  ));

  return (
    <div>
      <form
        className="bg-white p-8 rounded-xl relative shadow-lg w-md z-0"
        onSubmit={(e) => handleSubmit(e)}
      >
        <span className=" w-0 h-0 border-[10px] -z-10 -top-[20px] right-20 border-b-white border-t-transparent border-x-transparent  absolute"></span>
        <span className=" w-[20px] h-0 border-b-[1px]  -z-10 border-solid top-0 right-20 border-white absolute"></span>
        <h2 className="text-2xl font-medium mb-6 text-center">
          Đăng ký tài khoản
        </h2>
        {renderForm}
        <button
          type="submit"
          disabled={isPending}
          className="w-full h-10 mt-7 bg-[#FF7601] text-white flex justify-center items-center py-2 rounded-lg cursor-pointer relative group z-0"
        >
          <span>
            {isPending ? (
              <AiOutlineLoading3Quarters className="animate-spin " />
            ) : (
              "Đăng ký"
            )}
          </span>

          <span className="absolute top-0 bottom-0 left-0 w-0 bg-[#06923E] group-hover:w-full transition-all duration-200 -z-50" />
        </button>
      </form>
    </div>
  );
};

export default FormSignUp;
