import Header from "../components/NavBar";
import { Outlet } from "react-router-dom";

const MainLayout = () => {
  return (
    <>
      <Header />
      <div className=" bg-gray-100 h-auto w-full pb-10">
        <div className="h-auto w-full md:w-9/10 lg:w-17/25 m-auto px-2 pt-5">
          <Outlet />
        </div>
      </div>
    </>
  );
};

export default MainLayout;
