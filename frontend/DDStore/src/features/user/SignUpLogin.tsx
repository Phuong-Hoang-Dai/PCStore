import FormSignUp from "../auth/FormSignUp";
import FormSignIn from "../auth/FormSignIn";
import { useEffect, useRef, useState } from "react";

type Panel = "signIn" | "signUp" | null;

const SignUpLogin = () => {
  const [panel, setPanel] = useState<Panel>(null);
  // Wraps the buttons too, so clicking a button toggles instead of counting as "outside".
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!panel) {
      return;
    }
    const handleClickOutside = (e: MouseEvent) => {
      if (!ref.current?.contains(e.target as Node)) {
        setPanel(null);
      }
    };
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        setPanel(null);
      }
    };
    document.addEventListener("mousedown", handleClickOutside);
    document.addEventListener("keydown", handleEscape);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      document.removeEventListener("keydown", handleEscape);
    };
  }, [panel]);

  const toggle = (next: Exclude<Panel, null>) =>
    setPanel((prev) => (prev === next ? null : next));

  return (
    <div
      ref={ref}
      className="h-full flex flex-row justify-between items-center relative pr-4"
    >
      <button
        type="button"
        onClick={() => toggle("signUp")}
        aria-expanded={panel === "signUp"}
        className="bg-[#ff6d4d] relative h-1/2 px-5 rounded-4xl text-white cursor-pointer hover:bg-[#ff876c]"
      >
        Đăng ký
      </button>
      {panel === "signUp" && (
        <div className="absolute top-20 right-12 z-10">
          <FormSignUp onSwitch={() => setPanel("signIn")} />
        </div>
      )}
      <span className="border-r-1 h-1/3 ml-2"></span>
      <button
        type="button"
        onClick={() => toggle("signIn")}
        aria-expanded={panel === "signIn"}
        className="h-1/2 text-black font-medium cursor-pointer ml-2"
      >
        Đăng nhập
      </button>
      {panel === "signIn" && (
        <div className="absolute top-20 -right-2 z-10">
          <FormSignIn IsForUser={true} onSwitch={() => setPanel("signUp")} />
        </div>
      )}
    </div>
  );
};

export default SignUpLogin;
