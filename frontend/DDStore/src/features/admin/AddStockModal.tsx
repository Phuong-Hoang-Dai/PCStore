import { useState } from "react";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import Modal from "../../components/Modal";
import type { AdminProduct } from "./admin_model";

const AddStockModal = ({
  product,
  currentStock,
  onSubmit,
  onClose,
}: {
  product: AdminProduct;
  currentStock: number;
  onSubmit: (quantity: number) => Promise<void>;
  onClose: () => void;
}) => {
  const [quantity, setQuantity] = useState(10);
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsPending(true);
    setError("");
    onSubmit(quantity).catch((err: Error) => {
      setIsPending(false);
      setError(err.message);
    });
  };

  return (
    <Modal title="Nhập kho" onClose={onClose}>
      <form onSubmit={handleSubmit} className="flex flex-col gap-4 text-sm">
        <div>
          <p className="font-medium">{product.name}</p>
          <p className="text-gray-500">
            Tồn kho hiện tại: <span className="font-medium text-black">{currentStock}</span>
          </p>
        </div>
        <label className="flex flex-col gap-1">
          <span className="font-medium">Số lượng nhập thêm</span>
          <input
            required
            autoFocus
            type="number"
            min={1}
            step={1}
            value={quantity}
            onChange={(e) => setQuantity(Number(e.target.value))}
            className="w-full border border-gray-300 p-2 rounded-lg"
          />
        </label>
        <p className="text-gray-500">
          Sau khi nhập: <span className="font-medium text-black">{currentStock + (quantity || 0)}</span>
        </p>
        {error && <p className="text-red-500">{error}</p>}
        <div className="flex justify-end gap-3">
          <button
            type="button"
            onClick={onClose}
            className="px-5 h-10 rounded-4xl border border-gray-300 cursor-pointer hover:bg-gray-100"
          >
            Hủy
          </button>
          <button
            type="submit"
            disabled={isPending || quantity < 1}
            className="px-5 h-10 min-w-24 rounded-4xl bg-[#06923E] text-white cursor-pointer hover:bg-[#08a948] disabled:bg-green-300 flex justify-center items-center"
          >
            {isPending ? <AiOutlineLoading3Quarters className="animate-spin" /> : "Nhập kho"}
          </button>
        </div>
      </form>
    </Modal>
  );
};

export default AddStockModal;
