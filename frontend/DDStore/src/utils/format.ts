const vndFormatter = new Intl.NumberFormat("vi-VN", {
  style: "currency",
  currency: "VND",
});

export const formatVND = (amount: number) => vndFormatter.format(amount);

export const formatDate = (value: string) =>
  new Date(value).toLocaleString("vi-VN");
