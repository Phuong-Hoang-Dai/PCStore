import { useState } from "react";
import { AiOutlineLoading3Quarters } from "react-icons/ai";

import { SignIn, TryGetUserByCookie } from "./user_service";
import { useDispatch } from "react-redux";
import { SetUser } from "./userReducer";

const FormSignIn = ({ IsForUser: IsForUse }: { IsForUser: boolean }) => {
  const dispatch = useDispatch();
  const [form, setForm] = useState({
    name: "",
    password: "",
  });
  const [isPending, setIsPending] = useState(false);
  const [isError, setIsError] = useState(false);
  const labels = ["Họ tên", "Mật khẩu"];

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsPending(true);
    SignIn({
      id: 0,
      name: form.name,
      email: "",
      password: form.password,
    })
      .then(() => {
        setIsPending(false);
        return TryGetUserByCookie();
      })
      .then((data) => dispatch(SetUser(data)))
      .catch(() => {
        setIsPending(false);
        setIsError(true);
      });
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    key: string
  ) => {
    e.preventDefault();
    setForm((prev) => ({ ...prev, [key]: e.target.value }));
  };

  const keys = Object.keys(form);

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
        className="bg-white p-8 rounded-xl shadow-lg  w-md relative"
        onSubmit={(e) => handleSubmit(e)}
      >
        {IsForUse && (
          <>
            {" "}
            <span className=" w-0 h-0 border-[10px] -z-10 -top-[20px] right-5 border-b-white border-t-transparent border-x-transparent  absolute"></span>
            <span className=" w-[20px] h-0 border-b-[1px]  -z-10 border-solid top-0 right-20 border-white absolute"></span>
          </>
        )}

        <h2 className="text-2xl font-medium mb-6 text-center">
          Đăng nhập tài khoản
        </h2>
        {renderForm}
        {isError && (
          <p className="text-red-500 text-sm my-2">
            Tài khoản hoặc mật khẩu không hợp lệ
          </p>
        )}
        <button
          type="submit"
          disabled={isPending}
          className="w-full h-10 mt-7 bg-[#FF7601] text-white flex justify-center items-center py-2 rounded-lg cursor-pointer relative group z-0"
        >
          <span>
            {isPending ? (
              <AiOutlineLoading3Quarters className="animate-spin " />
            ) : (
              "Đăng nhập"
            )}
          </span>
          <span className="absolute top-0 bottom-0 left-0 w-0 bg-[#06923E] group-hover:w-full transition-all duration-200 -z-50" />
        </button>
      </form>
    </div>
  );
};

export default FormSignIn;
