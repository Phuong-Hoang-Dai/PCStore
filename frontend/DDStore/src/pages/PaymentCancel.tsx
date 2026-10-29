import { Link } from "react-router-dom";
import { FaRegCircleXmark } from "react-icons/fa6";

const PaymentCancel = () => {
  return (
    <div className="bg-white rounded-2xl shadow-lg max-w-lg mx-auto mt-10 p-10 flex flex-col items-center gap-4 text-center">
      <FaRegCircleXmark className="text-red-500" size={64} />
      <h1 className="text-2xl font-medium">Thanh toán đã bị hủy</h1>
      <p className="text-gray-600">
        Giỏ hàng của bạn vẫn được giữ lại. Bạn có thể thanh toán lại đơn hàng
        trong mục "Đơn hàng".
      </p>
      <div className="flex gap-3 mt-4">
        <Link
          to="/orders"
          className="px-6 py-2 rounded-4xl bg-[#FF7601] text-white hover:bg-[#ff8f33]"
        >
          Xem đơn hàng
        </Link>
        <Link
          to="/"
          className="px-6 py-2 rounded-4xl border border-gray-300 hover:bg-gray-100"
        >
          Về trang chủ
        </Link>
      </div>
    </div>
  );
};

export default PaymentCancel;
