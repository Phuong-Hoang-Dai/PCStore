import { useState } from "react";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import Modal from "../../components/Modal";
import { FALLBACK_IMAGE } from "../product/product_model";
import type { Category, ProductInput } from "./admin_model";

const emptyInput: ProductInput = {
  name: "",
  description: "",
  brand: "",
  categoryId: 0,
  price: 0,
  imageUrl: "",
};

const inputClass = "w-full border border-gray-300 p-2 rounded-lg";

const ProductFormModal = ({
  title,
  initial,
  categories,
  onSubmit,
  onClose,
}: {
  title: string;
  initial?: ProductInput;
  categories: Category[];
  onSubmit: (input: ProductInput) => Promise<void>;
  onClose: () => void;
}) => {
  const [form, setForm] = useState<ProductInput>(
    initial ?? { ...emptyInput, categoryId: categories[0]?.id ?? 0 }
  );
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState("");

  const set = <K extends keyof ProductInput>(key: K, value: ProductInput[K]) =>
    setForm((prev) => ({ ...prev, [key]: value }));

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    if (!form.categoryId) {
      setError("Vui lòng chọn danh mục");
      return;
    }
    setIsPending(true);
    setError("");
    onSubmit({ ...form, name: form.name.trim(), imageUrl: form.imageUrl.trim() }).catch(
      (err: Error) => {
        setIsPending(false);
        setError(err.message);
      }
    );
  };

  return (
    <Modal title={title} onClose={onClose}>
      <form onSubmit={handleSubmit} className="flex flex-col gap-4 text-sm">
        <label className="flex flex-col gap-1">
          <span className="font-medium">Tên sản phẩm</span>
          <input
            required
            value={form.name}
            onChange={(e) => set("name", e.target.value)}
            className={inputClass}
          />
        </label>
        <label className="flex flex-col gap-1">
          <span className="font-medium">Mô tả</span>
          <textarea
            rows={3}
            value={form.description}
            onChange={(e) => set("description", e.target.value)}
            className={inputClass}
          />
        </label>
        <div className="grid grid-cols-2 gap-4">
          <label className="flex flex-col gap-1">
            <span className="font-medium">Thương hiệu</span>
            <input
              value={form.brand}
              onChange={(e) => set("brand", e.target.value)}
              className={inputClass}
            />
          </label>
          <label className="flex flex-col gap-1">
            <span className="font-medium">Danh mục</span>
            <select
              required
              value={form.categoryId || ""}
              onChange={(e) => set("categoryId", Number(e.target.value))}
              className={inputClass}
            >
              <option value="" disabled>
                Chọn danh mục
              </option>
              {categories.map((c) => (
                <option key={c.id} value={c.id}>
                  {c.name}
                </option>
              ))}
            </select>
          </label>
        </div>
        <label className="flex flex-col gap-1">
          <span className="font-medium">Giá (VND)</span>
          <input
            required
            type="number"
            min={0}
            step={1000}
            value={form.price}
            onChange={(e) => set("price", Number(e.target.value))}
            className={inputClass}
          />
        </label>
        <label className="flex flex-col gap-1">
          <span className="font-medium">Ảnh (URL)</span>
          <div className="flex gap-3 items-center">
            <input
              type="url"
              placeholder="https://..."
              value={form.imageUrl}
              onChange={(e) => set("imageUrl", e.target.value)}
              className={inputClass}
            />
            <img
              src={form.imageUrl || FALLBACK_IMAGE}
              alt=""
              className="w-12 h-12 rounded-lg object-cover border border-gray-200 shrink-0"
            />
          </div>
        </label>
        {error && <p className="text-red-500">{error}</p>}
        <div className="flex justify-end gap-3 mt-2">
          <button
            type="button"
            onClick={onClose}
            className="px-5 h-10 rounded-4xl border border-gray-300 cursor-pointer hover:bg-gray-100"
          >
            Hủy
          </button>
          <button
            type="submit"
            disabled={isPending}
            className="px-5 h-10 min-w-24 rounded-4xl bg-[#FF7601] text-white cursor-pointer hover:bg-[#ff8f33] disabled:bg-orange-300 flex justify-center items-center"
          >
            {isPending ? <AiOutlineLoading3Quarters className="animate-spin" /> : "Lưu"}
          </button>
        </div>
      </form>
    </Modal>
  );
};

export default ProductFormModal;
