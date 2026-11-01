import { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import { selectUser } from "../features/auth/userReducer";
import {
  canPay,
  ORDER_PAYMENT_TIMEOUT_HOURS,
  orderStateColors,
  orderStateLabels,
  type Order,
} from "../features/order/order_model";
import {
  createCheckoutSession,
  fetchOrderHistory,
} from "../features/order/order_service";
import { formatDate, formatVND } from "../utils/format";

const OrderHistory = () => {
  const user = useSelector(selectUser);
  const [orders, setOrders] = useState<Order[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [payingId, setPayingId] = useState<number | null>(null);

  const isLoggedIn = user.id !== 0;

  useEffect(() => {
    if (!isLoggedIn) {
      return;
    }
    fetchOrderHistory()
      .then((data) => setOrders(data))
      .catch((err: Error) => setError(err.message))
      .finally(() => setIsLoading(false));
  }, [isLoggedIn]);

  const handlePay = (orderId: number) => {
    setPayingId(orderId);
    setError("");
    createCheckoutSession(orderId)
      .then((url) => window.location.assign(url))
      .catch((err: Error) => {
        setPayingId(null);
        setError(err.message);
      });
  };

  if (!isLoggedIn) {
    return (
      <div className="bg-white rounded-2xl shadow-lg p-10 text-center text-gray-600">
        Vui lòng đăng nhập để xem đơn hàng của bạn.
      </div>
    );
  }

  const renderOrders = orders.map((order) => (
    <div
      key={order.orderId}
      className="bg-white rounded-2xl shadow p-5 flex flex-col md:flex-row md:items-center justify-between gap-3"
    >
      <div className="flex flex-col gap-1">
        <span className="font-medium">Đơn hàng #{order.orderId}</span>
        <span className="text-sm text-gray-500">
          {formatDate(order.create_at)}
        </span>
      </div>
      <div className="flex flex-row items-center gap-4">
        <span
          className={`text-xs font-medium px-3 py-1 rounded-4xl ${
            orderStateColors[order.state] ?? "bg-gray-100"
          }`}
        >
          {orderStateLabels[order.state] ?? "Không xác định"}
        </span>
        <span className="text-red-600 font-medium min-w-28 text-end">
          {formatVND(order.total)}
        </span>
        {canPay(order.state) && (
          <button
            type="button"
            disabled={payingId !== null}
            onClick={() => handlePay(order.orderId)}
            className="h-9 w-28 rounded-4xl bg-[#FF7601] text-white cursor-pointer hover:bg-[#ff8f33] disabled:bg-orange-300 disabled:cursor-not-allowed flex justify-center items-center"
          >
            {payingId === order.orderId ? (
              <AiOutlineLoading3Quarters className="animate-spin" />
            ) : (
              "Thanh toán"
            )}
          </button>
        )}
      </div>
    </div>
  ));

  return (
    <div className="flex flex-col gap-4">
      <h1 className="text-2xl font-medium">Đơn hàng của tôi</h1>
      <p className="text-sm text-gray-500">
        Đơn hàng chưa thanh toán sẽ tự động bị hủy sau {ORDER_PAYMENT_TIMEOUT_HOURS} giờ.
      </p>
      {error && <p className="text-red-500 text-sm">{error}</p>}
      {isLoading ? (
        <div className="flex justify-center p-10">
          <AiOutlineLoading3Quarters className="animate-spin" size={32} />
        </div>
      ) : orders.length === 0 ? (
        <div className="bg-white rounded-2xl shadow p-10 text-center text-gray-500">
          Bạn chưa có đơn hàng nào.
        </div>
      ) : (
        renderOrders
      )}
    </div>
  );
};

export default OrderHistory;
