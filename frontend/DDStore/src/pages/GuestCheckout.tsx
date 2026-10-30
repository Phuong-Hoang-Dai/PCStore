import { useState } from "react";
import { useSelector } from "react-redux";
import { Link } from "react-router-dom";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import { selectCart } from "../features/order/cartSlice";
import { selectUser } from "../features/auth/userReducer";
import { createGuestOrder } from "../features/order/order_service";
import { ORDER_PAYMENT_TIMEOUT_HOURS } from "../features/order/order_model";
import { formatVND } from "../utils/format";

// Checkout for customers who aren't logged in: they only give a name and an email.
const GuestCheckout = () => {
  const cart = useSelector(selectCart);
  const user = useSelector(selectUser);
  const [form, setForm] = useState({ name: "", email: "" });
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState("");

  const total = cart.order.reduce(
    (sum, item) => sum + item.price * item.quantityOrder,
    0
  );

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (isPending) {
      return;
    }
    setIsPending(true);
    setError("");
    createGuestOrder({
      name: form.name.trim(),
      email: form.email.trim(),
      data: cart.order.map((item) => ({
        productId: item.id,
        quantity: item.quantityOrder,
      })),
    })
      // The cart is cleared on the success page, so it survives a cancelled payment.
      .then((url) => window.location.assign(url))
      .catch((err: Error) => {
        setIsPending(false);
        setError(err.message || "Không thể tạo đơn hàng, vui lòng thử lại");
      });
  };

  if (cart.order.length === 0) {
    return (
      <div className="bg-white rounded-2xl shadow-lg max-w-lg mx-auto mt-10 p-10 text-center text-gray-600">
        <p>Giỏ hàng trống.</p>
        <Link to="/" className="inline-block mt-4 text-[#FF7601] hover:underline">
          Tiếp tục mua sắm
        </Link>
      </div>
    );
  }

  if (user.id !== 0) {
    return (
      <div className="bg-white rounded-2xl shadow-lg max-w-lg mx-auto mt-10 p-10 text-center text-gray-600">
        Bạn đã đăng nhập. Hãy mở giỏ hàng và chọn "Hoàn tất đơn hàng".
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto mt-6 grid gap-6 md:grid-cols-2">
      <form
        onSubmit={handleSubmit}
        className="bg-white rounded-2xl shadow-lg p-8 flex flex-col gap-4"
      >
        <h1 className="text-2xl font-medium">Thông tin đặt hàng</h1>
        <p className="text-sm text-gray-600">
          Bạn không cần tạo tài khoản. Nếu sau này đăng ký bằng email này, các
          đơn hàng sẽ hiện trong lịch sử đơn hàng của bạn.
        </p>
        <label className="flex flex-col gap-1">
          <span className="font-medium">Họ tên</span>
          <input
            required
            autoComplete="name"
            value={form.name}
            onChange={(e) => setForm((prev) => ({ ...prev, name: e.target.value }))}
            className="border border-gray-300 p-2 rounded-lg"
          />
        </label>
        <label className="flex flex-col gap-1">
          <span className="font-medium">Email</span>
          <input
            required
            type="email"
            autoComplete="email"
            value={form.email}
            onChange={(e) => setForm((prev) => ({ ...prev, email: e.target.value }))}
            className="border border-gray-300 p-2 rounded-lg"
          />
        </label>
        {error && <p className="text-red-500 text-sm">{error}</p>}
        <button
          type="submit"
          disabled={isPending}
          className="h-12 mt-2 bg-red-600 text-white text-lg rounded-4xl hover:bg-red-400 disabled:bg-red-300 disabled:cursor-not-allowed flex justify-center items-center cursor-pointer"
        >
          {isPending ? (
            <AiOutlineLoading3Quarters className="animate-spin" />
          ) : (
            "Đặt hàng và thanh toán"
          )}
        </button>
        <p className="text-xs text-gray-500">
          Đơn hàng chưa thanh toán sẽ tự động hủy sau{" "}
          {ORDER_PAYMENT_TIMEOUT_HOURS} giờ.
        </p>
      </form>

      <div className="bg-white rounded-2xl shadow-lg p-8 flex flex-col gap-3">
        <h2 className="text-xl font-medium">Đơn hàng</h2>
        <ul className="flex flex-col gap-3">
          {cart.order.map((item) => (
            <li key={item.id} className="flex items-center gap-3 text-sm">
              <img
                src={item.image}
                alt={item.name}
                className="w-14 aspect-square object-cover rounded"
              />
              <span className="flex-1 truncate">
                {item.name} × {item.quantityOrder}
              </span>
              <span>{formatVND(item.price * item.quantityOrder)}</span>
            </li>
          ))}
        </ul>
        <div className="flex justify-between border-t border-gray-200 pt-3 text-lg font-medium">
          <span>Tổng tiền</span>
          <span>{formatVND(total)}</span>
        </div>
      </div>
    </div>
  );
};

export default GuestCheckout;
