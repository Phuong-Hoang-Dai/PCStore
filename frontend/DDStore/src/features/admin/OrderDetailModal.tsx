import { useEffect, useState } from "react";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import Modal from "../../components/Modal";
import { orderStateColors, orderStateLabels, type OrderDetail } from "../order/order_model";
import { formatDate, formatVND } from "../../utils/format";
import { fetchOrderDetail, fetchProduct } from "./admin_service";

const OrderDetailModal = ({ orderId, onClose }: { orderId: number; onClose: () => void }) => {
  const [order, setOrder] = useState<OrderDetail | null>(null);
  const [names, setNames] = useState<Record<string, string>>({});
  const [error, setError] = useState("");

  useEffect(() => {
    let ignore = false;
    fetchOrderDetail(orderId)
      .then(async (data) => {
        if (ignore) {
          return;
        }
        setOrder(data);
        // Orders only store product ids; a deleted product just keeps showing its id.
        const ids = [...new Set((data.items ?? []).map((i) => i.productId))];
        const results = await Promise.allSettled(ids.map(fetchProduct));
        const found: Record<string, string> = {};
        results.forEach((r, i) => {
          if (r.status === "fulfilled") {
            found[ids[i]] = r.value.name;
          }
        });
        if (!ignore) {
          setNames(found);
        }
      })
      .catch((err: Error) => !ignore && setError(err.message));
    return () => {
      ignore = true;
    };
  }, [orderId]);

  return (
    <Modal title={`Đơn hàng #${orderId}`} onClose={onClose}>
      {error ? (
        <p className="text-red-500 text-sm">{error}</p>
      ) : !order ? (
        <div className="flex justify-center p-6">
          <AiOutlineLoading3Quarters className="animate-spin" size={28} />
        </div>
      ) : (
        <div className="flex flex-col gap-4 text-sm">
          <div className="flex flex-wrap gap-x-6 gap-y-1 text-gray-600">
            <span>Khách hàng: #{order.userId}</span>
            <span>{formatDate(order.create_at)}</span>
            <span
              className={`text-xs font-medium px-3 py-0.5 rounded-4xl ${
                orderStateColors[order.state] ?? "bg-gray-100"
              }`}
            >
              {orderStateLabels[order.state] ?? "Không xác định"}
            </span>
          </div>
          <table className="w-full">
            <thead>
              <tr className="text-left text-gray-500">
                <th className="pb-2 font-medium">Sản phẩm</th>
                <th className="pb-2 font-medium text-right">SL</th>
                <th className="pb-2 font-medium text-right">Thành tiền</th>
              </tr>
            </thead>
            <tbody>
              {(order.items ?? []).map((item) => (
                <tr key={item.productId} className="border-t border-gray-100">
                  <td className="py-2 pr-3">
                    <div className="font-medium">{names[item.productId] ?? item.productId}</div>
                    <div className="text-gray-500">{formatVND(item.price)}</div>
                  </td>
                  <td className="py-2 text-right">{item.quantity}</td>
                  <td className="py-2 text-right">{formatVND(item.price * item.quantity)}</td>
                </tr>
              ))}
            </tbody>
          </table>
          <div className="flex justify-between border-t border-gray-200 pt-3 text-base font-medium">
            <span>Tổng tiền</span>
            <span className="text-red-600">{formatVND(order.total)}</span>
          </div>
        </div>
      )}
    </Modal>
  );
};

export default OrderDetailModal;
