import { useDispatch, useSelector } from "react-redux";
import { Link, useSearchParams } from "react-router-dom";
import { RiLogoutBoxRLine } from "react-icons/ri";
import { RemoveUser, selectUser } from "../features/auth/userReducer";
import { Logout } from "../features/auth/user_service";
import ProductsTab from "../features/admin/ProductsTab";
import CategoriesTab from "../features/admin/CategoriesTab";
import OrdersTab from "../features/admin/OrdersTab";

const tabs = [
  { key: "products", label: "Sản phẩm", element: <ProductsTab /> },
  { key: "categories", label: "Danh mục", element: <CategoriesTab /> },
  { key: "orders", label: "Đơn hàng", element: <OrdersTab /> },
] as const;

const Manage = () => {
  const user = useSelector(selectUser);
  const dispatch = useDispatch();
  // The active tab lives in the URL so a refresh or shared link opens the same tab.
  const [searchParams, setSearchParams] = useSearchParams();
  const active = tabs.find((t) => t.key === searchParams.get("tab")) ?? tabs[0];

  const handleLogout = () => {
    // RequreAdmin redirects to the admin login once the user is cleared.
    Logout().finally(() => dispatch(RemoveUser()));
  };

  return (
    <div className="min-h-screen bg-gray-100">
      <header className="bg-[#ecb996] h-16 w-full">
        <div className="flex items-center justify-between h-full w-full md:w-9/10 lg:w-17/25 m-auto px-4">
          <div className="flex items-center gap-3 h-full py-2">
            <Link to="/" className="h-full">
              <img src="/logo.png" alt="DDStore" className="h-full" />
            </Link>
            <span className="font-medium text-lg">Quản trị</span>
          </div>
          <div className="flex items-center gap-4 text-sm font-bold">
            <Link to="/" className="hover:underline hidden md:inline">
              Về cửa hàng
            </Link>
            <span>
              {user.name} <span className="font-normal">({user.roleId})</span>
            </span>
            <button
              type="button"
              onClick={handleLogout}
              aria-label="Đăng xuất"
              className="cursor-pointer"
            >
              <RiLogoutBoxRLine size={18} />
            </button>
          </div>
        </div>
      </header>

      <main className="w-full md:w-9/10 lg:w-17/25 m-auto px-4 py-6 flex flex-col gap-5">
        <nav className="flex gap-2 border-b border-gray-300" role="tablist">
          {tabs.map((tab) => (
            <button
              key={tab.key}
              type="button"
              role="tab"
              aria-selected={tab.key === active.key}
              onClick={() => setSearchParams({ tab: tab.key })}
              className={`px-4 py-2 -mb-px border-b-2 cursor-pointer font-medium ${
                tab.key === active.key
                  ? "border-[#FF7601] text-[#FF7601]"
                  : "border-transparent text-gray-600 hover:text-black"
              }`}
            >
              {tab.label}
            </button>
          ))}
        </nav>
        {active.element}
      </main>
    </div>
  );
};

export default Manage;
