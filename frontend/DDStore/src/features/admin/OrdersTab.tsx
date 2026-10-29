import { useEffect, useState } from "react";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import {
  canPay,
  ORDER_PAYMENT_TIMEOUT_HOURS,
  OrderState,
  orderStateColors,
  orderStateLabels,
  type Order,
} from "../order/order_model";
import { formatDate, formatVND } from "../../utils/format";
import { cancelOrder, completeOrder, fetchOrders } from "./admin_service";
import OrderDetailModal from "./OrderDetailModal";
import Pager from "./Pager";

const PAGE_SIZE = 20;

const OrdersTab = () => {
  const [orders, setOrders] = useState<Order[]>([]);
  const [page, setPage] = useState(0);
  const [reload, setReload] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [busyId, setBusyId] = useState<number | null>(null);
  const [viewingId, setViewingId] = useState<number | null>(null);

  useEffect(() => {
    let ignore = false;
    setIsLoading(true);
    fetchOrders(PAGE_SIZE, page * PAGE_SIZE)
      .then((data) => {
        if (!ignore) {
          setOrders(data);
          setError("");
        }
      })
      .catch((err: Error) => !ignore && setError(err.message))
      .finally(() => !ignore && setIsLoading(false));
    return () => {
      ignore = true;
    };
  }, [page, reload]);

  const run = (orderId: number, action: (id: number) => Promise<unknown>, confirmText: string) => {
    if (!window.confirm(confirmText)) {
      return;
    }
    setBusyId(orderId);
    action(orderId)
      .then(() => setReload((n) => n + 1))
      .catch((err: Error) => setError(err.message))
      .finally(() => setBusyId(null));
  };

  const rows = orders.map((order) => (
    <tr key={order.orderId} className="border-t border-gray-100">
      <td className="py-3 pr-3 font-medium">#{order.orderId}</td>
      <td className="py-3 pr-3 text-gray-600">#{order.userId}</td>
      <td className="py-3 pr-3 whitespace-nowrap text-gray-600">{formatDate(order.create_at)}</td>
      <td className="py-3 pr-3 whitespace-nowrap text-right">{formatVND(order.total)}</td>
      <td className="py-3 pr-3">
        <span
          className={`text-xs font-medium px-3 py-1 rounded-4xl whitespace-nowrap ${
            orderStateColors[order.state] ?? "bg-gray-100"
          }`}
        >
          {orderStateLabels[order.state] ?? "Không xác định"}
        </span>
      </td>
      <td className="py-3">
        <div className="flex gap-3 justify-end whitespace-nowrap">
          {busyId === order.orderId ? (
            <AiOutlineLoading3Quarters className="animate-spin" />
          ) : (
            <>
              <button
                type="button"
                onClick={() => setViewingId(order.orderId)}
                className="text-gray-700 font-medium cursor-pointer hover:underline"
              >
                Chi tiết
              </button>
              {order.state === OrderState.Payed && (
                <button
                  type="button"
                  onClick={() =>
                    run(order.orderId, completeOrder, `Hoàn tất đơn hàng #${order.orderId}?`)
                  }
                  className="text-[#06923E] font-medium cursor-pointer hover:underline"
                >
                  Hoàn tất
                </button>
              )}
              {canPay(order.state) && (
                <button
                  type="button"
                  onClick={() =>
                    run(
                      order.orderId,
                      cancelOrder,
                      `Hủy đơn hàng #${order.orderId}? Hàng sẽ được trả lại kho.`
                    )
                  }
                  className="text-red-600 font-medium cursor-pointer hover:underline"
                >
                  Hủy
                </button>
              )}
            </>
          )}
        </div>
      </td>
    </tr>
  ));

  return (
    <div className="flex flex-col gap-4">
      <p className="text-xs text-gray-500">
        Đơn chưa thanh toán sẽ tự động bị hủy và trả hàng về kho sau{" "}
        {ORDER_PAYMENT_TIMEOUT_HOURS} giờ.
      </p>
      {error && <p className="bg-red-50 text-red-700 text-sm rounded-lg px-4 py-2">{error}</p>}

      <div className="bg-white rounded-2xl shadow p-5 overflow-x-auto">
        {isLoading ? (
          <div className="flex justify-center p-10">
            <AiOutlineLoading3Quarters className="animate-spin" size={32} />
          </div>
        ) : orders.length === 0 ? (
          <p className="text-center text-gray-500 p-10">Chưa có đơn hàng nào.</p>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-gray-500">
                <th className="pb-2 font-medium">Mã đơn</th>
                <th className="pb-2 font-medium">Khách hàng</th>
                <th className="pb-2 font-medium">Ngày tạo</th>
                <th className="pb-2 font-medium text-right pr-3">Tổng tiền</th>
                <th className="pb-2 font-medium">Trạng thái</th>
                <th className="pb-2" />
              </tr>
            </thead>
            <tbody>{rows}</tbody>
          </table>
        )}
        <Pager page={page} hasNext={orders.length === PAGE_SIZE} onChange={setPage} />
      </div>

      {viewingId !== null && (
        <OrderDetailModal orderId={viewingId} onClose={() => setViewingId(null)} />
      )}
    </div>
  );
};

export default OrdersTab;
