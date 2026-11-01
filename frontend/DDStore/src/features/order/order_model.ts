
// Mirrors model.OrderState in order_service.
export const OrderState = {
  Processing: 0,
  ReadyToPay: 1,
  Payed: 2,
  Completed: 3,
  Cancelled: 4,
} as const;

export const orderStateLabels: Record<number, string> = {
  [OrderState.Processing]: "Đang xử lý",
  [OrderState.ReadyToPay]: "Chờ thanh toán",
  [OrderState.Payed]: "Đã thanh toán",
  [OrderState.Completed]: "Hoàn tất",
  [OrderState.Cancelled]: "Đã hủy",
};

export const orderStateColors: Record<number, string> = {
  [OrderState.Processing]: "bg-yellow-100 text-yellow-800",
  [OrderState.ReadyToPay]: "bg-orange-100 text-orange-800",
  [OrderState.Payed]: "bg-green-100 text-green-800",
  [OrderState.Completed]: "bg-blue-100 text-blue-800",
  [OrderState.Cancelled]: "bg-gray-200 text-gray-600",
};

export const canPay = (state: number) =>
  state === OrderState.Processing || state === OrderState.ReadyToPay;

// Unpaid orders are cancelled by order_service this long after they are created.
export const ORDER_PAYMENT_TIMEOUT_HOURS = 3;

export type Order = {
  orderId: number;
  state: number;
  total: number;
  userId: number;
  create_at: string;
}

export type OrderItem = {
  productId: string;
  quantity: number;
  price: number;
}

export type OrderDetail = Order & {
  items: OrderItem[] | null;
}

export type OrderItemRequest = {
  productId: string;
  quantity: number;
}
