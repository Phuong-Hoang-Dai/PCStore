import { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import { selectUser } from "../auth/userReducer";
import { FALLBACK_IMAGE, ProductState } from "../product/product_model";
import { formatVND } from "../../utils/format";
import { toProductInput, type AdminProduct, type Category } from "./admin_model";
import {
  addStock,
  createProduct,
  deleteProduct,
  fetchCategories,
  fetchProducts,
  fetchStocks,
  updateProduct,
} from "./admin_service";
import ProductFormModal from "./ProductFormModal";
import AddStockModal from "./AddStockModal";
import Pager from "./Pager";

const PAGE_SIZE = 10;
// The gateway queues new products; the catalog usually creates them well within this.
const CREATE_REFRESH_DELAY_MS = 1500;

const ProductsTab = () => {
  const user = useSelector(selectUser);
  const [products, setProducts] = useState<AdminProduct[]>([]);
  const [stocks, setStocks] = useState<Record<string, number>>({});
  const [categories, setCategories] = useState<Category[]>([]);
  const [categoryId, setCategoryId] = useState(0);
  const [page, setPage] = useState(0);
  const [reload, setReload] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  // null: closed, "new": create form, otherwise the product being edited.
  const [editing, setEditing] = useState<AdminProduct | "new" | null>(null);
  const [stocking, setStocking] = useState<AdminProduct | null>(null);

  const refresh = () => setReload((n) => n + 1);

  useEffect(() => {
    fetchCategories()
      .then(setCategories)
      .catch((err: Error) => setError(err.message));
  }, []);

  useEffect(() => {
    let ignore = false;
    setIsLoading(true);
    fetchProducts(PAGE_SIZE, page * PAGE_SIZE, categoryId || undefined)
      .then(async (data) => {
        const stockById = await fetchStocks(data.map((p) => p.id));
        if (!ignore) {
          setProducts(data);
          setStocks(stockById);
          setError("");
        }
      })
      .catch((err: Error) => !ignore && setError(err.message))
      .finally(() => !ignore && setIsLoading(false));
    return () => {
      ignore = true;
    };
  }, [page, categoryId, reload]);

  const handleCreate = async (input: Parameters<typeof createProduct>[0]) => {
    await createProduct(input);
    setEditing(null);
    setNotice(`Đang tạo "${input.name}". Sản phẩm mới ở trạng thái hết hàng cho đến khi nhập kho.`);
    setTimeout(refresh, CREATE_REFRESH_DELAY_MS);
  };

  const handleUpdate = async (id: string, input: Parameters<typeof updateProduct>[1]) => {
    await updateProduct(id, input);
    setEditing(null);
    setNotice(`Đã cập nhật "${input.name}".`);
    refresh();
  };

  const handleAddStock = async (product: AdminProduct, quantity: number) => {
    const total = await addStock(product.id, quantity);
    setStocking(null);
    setNotice(`Đã nhập ${quantity} "${product.name}". Tồn kho: ${total}.`);
    refresh();
  };

  const handleDelete = (product: AdminProduct) => {
    if (!window.confirm(`Xóa sản phẩm "${product.name}"?`)) {
      return;
    }
    deleteProduct(product.id)
      .then(() => {
        setNotice(`Đã xóa "${product.name}".`);
        refresh();
      })
      .catch((err: Error) => setError(err.message));
  };

  const rows = products.map((product) => {
    const stock = stocks[product.id] ?? 0;
    const isAvailable = product.state === ProductState.IsAvailable;
    return (
      <tr key={product.id} className="border-t border-gray-100 align-middle">
        <td className="py-3 pr-3">
          <img
            src={product.images?.[0]?.url || FALLBACK_IMAGE}
            alt=""
            className="w-12 h-12 rounded-lg object-cover"
          />
        </td>
        <td className="py-3 pr-3 max-w-64">
          <div className="font-medium truncate">{product.name}</div>
          <div className="text-gray-500 truncate">{product.description}</div>
        </td>
        <td className="py-3 pr-3 whitespace-nowrap">{product.cate?.name || "—"}</td>
        <td className="py-3 pr-3 whitespace-nowrap text-right">{formatVND(product.price)}</td>
        <td className={`py-3 pr-3 text-right font-medium ${stock === 0 ? "text-red-600" : ""}`}>
          {stock}
        </td>
        <td className="py-3 pr-3">
          <span
            className={`text-xs font-medium px-3 py-1 rounded-4xl whitespace-nowrap ${
              isAvailable ? "bg-green-100 text-green-800" : "bg-gray-200 text-gray-600"
            }`}
          >
            {isAvailable ? "Còn hàng" : "Hết hàng"}
          </span>
        </td>
        <td className="py-3">
          <div className="flex gap-3 justify-end whitespace-nowrap">
            <button
              type="button"
              onClick={() => setStocking(product)}
              className="text-[#06923E] font-medium cursor-pointer hover:underline"
            >
              Nhập kho
            </button>
            <button
              type="button"
              onClick={() => setEditing(product)}
              className="text-[#FF7601] font-medium cursor-pointer hover:underline"
            >
              Sửa
            </button>
            {/* The gateway only lets the owner delete products. */}
            {user.roleId === "owner" && (
              <button
                type="button"
                onClick={() => handleDelete(product)}
                className="text-red-600 font-medium cursor-pointer hover:underline"
              >
                Xóa
              </button>
            )}
          </div>
        </td>
      </tr>
    );
  });

  return (
    <div className="flex flex-col gap-4">
      <div className="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <select
          value={categoryId}
          onChange={(e) => {
            setCategoryId(Number(e.target.value));
            setPage(0);
          }}
          className="border border-gray-300 p-2 rounded-lg bg-white text-sm md:w-56"
          aria-label="Lọc theo danh mục"
        >
          <option value={0}>Tất cả danh mục</option>
          {categories.map((c) => (
            <option key={c.id} value={c.id}>
              {c.name}
            </option>
          ))}
        </select>
        <button
          type="button"
          onClick={() => setEditing("new")}
          disabled={categories.length === 0}
          title={categories.length === 0 ? "Hãy tạo danh mục trước" : undefined}
          className="px-5 h-10 rounded-4xl bg-[#FF7601] text-white cursor-pointer hover:bg-[#ff8f33] disabled:bg-orange-300 disabled:cursor-not-allowed"
        >
          + Thêm sản phẩm
        </button>
      </div>

      {notice && (
        <p className="bg-green-50 text-green-800 text-sm rounded-lg px-4 py-2">{notice}</p>
      )}
      {error && <p className="bg-red-50 text-red-700 text-sm rounded-lg px-4 py-2">{error}</p>}

      <div className="bg-white rounded-2xl shadow p-5 overflow-x-auto">
        {isLoading ? (
          <div className="flex justify-center p-10">
            <AiOutlineLoading3Quarters className="animate-spin" size={32} />
          </div>
        ) : products.length === 0 ? (
          <p className="text-center text-gray-500 p-10">Chưa có sản phẩm nào.</p>
        ) : (
          <table className="w-full text-sm">
            <thead>
              <tr className="text-left text-gray-500">
                <th className="pb-2 font-medium" />
                <th className="pb-2 font-medium">Sản phẩm</th>
                <th className="pb-2 font-medium">Danh mục</th>
                <th className="pb-2 font-medium text-right pr-3">Giá</th>
                <th className="pb-2 font-medium text-right pr-3">Tồn kho</th>
                <th className="pb-2 font-medium">Trạng thái</th>
                <th className="pb-2" />
              </tr>
            </thead>
            <tbody>{rows}</tbody>
          </table>
        )}
        <Pager page={page} hasNext={products.length === PAGE_SIZE} onChange={setPage} />
      </div>

      {editing === "new" && (
        <ProductFormModal
          title="Thêm sản phẩm"
          categories={categories}
          onSubmit={handleCreate}
          onClose={() => setEditing(null)}
        />
      )}
      {editing && editing !== "new" && (
        <ProductFormModal
          title="Sửa sản phẩm"
          initial={toProductInput(editing)}
          categories={categories}
          onSubmit={(input) => handleUpdate(editing.id, input)}
          onClose={() => setEditing(null)}
        />
      )}
      {stocking && (
        <AddStockModal
          product={stocking}
          currentStock={stocks[stocking.id] ?? 0}
          onSubmit={(quantity) => handleAddStock(stocking, quantity)}
          onClose={() => setStocking(null)}
        />
      )}
    </div>
  );
};

export default ProductsTab;
