import FormSignUp from "../auth/FormSignUp";
import FormSignIn from "../auth/FormSignIn";
import { useEffect, useRef, useState } from "react";

const SignUpLogin = () => {
  const [isSignIn, setIsSignIn] = useState(false);
  const refSigIn = useRef<HTMLDivElement>(null);
  const [isSignUp, setIsSignUp] = useState(false);
  const refSigUp = useRef<HTMLDivElement>(null);
  const handleOnSignUp = (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) => {
    e.preventDefault();
    setIsSignUp(true);
    setIsSignIn(false);
  };

  const handleOnSignIn = (
    e: React.MouseEvent<HTMLButtonElement, MouseEvent>
  ) => {
    e.preventDefault();
    setIsSignUp(false);
    setIsSignIn(true);
  };
  const handleClickOutside = (e: MouseEvent) => {
    if (!refSigUp.current?.contains(e.target as Node)) {
      setIsSignUp(false);
    }
    if (!refSigIn.current?.contains(e.target as Node)) {
      setIsSignIn(false);
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
      <div className="h-full flex flex-row justify-between items-center relative pr-4">
        <button
          onClick={(e) => handleOnSignUp(e)}
          className="bg-[#ff6d4d] relative h-1/2 px-5 rounded-4xl text-white cursor-pointer hover:bg-[#ff876c]"
        >
          Sign up
        </button>
        {isSignUp && (
          <div className="absolute top-20 right-12 z-10" ref={refSigUp}>
            <FormSignUp />
          </div>
        )}
        <span className="border-r-1 h-1/3 ml-2"></span>
        <button
          onClick={(e) => handleOnSignIn(e)}
          className="h-1/2 text-black font-medium cursor-pointer ml-2"
        >
          Sign in
        </button>
        {isSignIn && (
          <div className="absolute top-20 -right-2 z-10" ref={refSigIn}>
            <FormSignIn IsForUser={true} />
          </div>
        )}
      </div>
    </>
  );
};

export default SignUpLogin;
