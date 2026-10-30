import { useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { Link } from "react-router-dom";
import { FaRegCircleCheck } from "react-icons/fa6";
import { clearCart } from "../features/order/cartSlice";
import { selectUser } from "../features/auth/userReducer";

const PaymentSuccess = () => {
  const dispatch = useDispatch();
  const isLoggedIn = useSelector(selectUser).id !== 0;

  useEffect(() => {
    dispatch(clearCart());
  }, [dispatch]);

  return (
    <div className="bg-white rounded-2xl shadow-lg max-w-lg mx-auto mt-10 p-10 flex flex-col items-center gap-4 text-center">
      <FaRegCircleCheck className="text-green-600" size={64} />
      <h1 className="text-2xl font-medium">Thanh toán thành công</h1>
      <p className="text-gray-600">
        Cảm ơn bạn đã mua hàng. Đơn hàng sẽ được cập nhật trạng thái "Đã thanh
        toán" trong giây lát.
      </p>
      {!isLoggedIn && (
        <p className="text-gray-600 text-sm">
          Đăng ký tài khoản bằng email bạn vừa dùng để xem lại đơn hàng này.
        </p>
      )}
      <div className="flex gap-3 mt-4">
        {isLoggedIn && (
          <Link
            to="/orders"
            className="px-6 py-2 rounded-4xl bg-[#FF7601] text-white hover:bg-[#ff8f33]"
          >
            Xem đơn hàng
          </Link>
        )}
        <Link
          to="/"
          className="px-6 py-2 rounded-4xl border border-gray-300 hover:bg-gray-100"
        >
          Tiếp tục mua sắm
        </Link>
      </div>
    </div>
  );
};

export default PaymentSuccess;
