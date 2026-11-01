import { useEffect, useState } from "react";
import { useSelector } from "react-redux";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import { selectUser } from "../auth/userReducer";
import type { Category } from "./admin_model";
import { createCategory, deleteCategory, fetchCategories, renameCategory } from "./admin_service";

const CategoriesTab = () => {
  const user = useSelector(selectUser);
  const [categories, setCategories] = useState<Category[]>([]);
  const [reload, setReload] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState("");
  const [newName, setNewName] = useState("");
  const [isAdding, setIsAdding] = useState(false);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editName, setEditName] = useState("");

  const refresh = () => setReload((n) => n + 1);

  useEffect(() => {
    let ignore = false;
    fetchCategories()
      .then((data) => !ignore && setCategories(data))
      .catch((err: Error) => !ignore && setError(err.message))
      .finally(() => !ignore && setIsLoading(false));
    return () => {
      ignore = true;
    };
  }, [reload]);

  const run = (action: Promise<unknown>) =>
    action
      .then(() => {
        setError("");
        refresh();
      })
      .catch((err: Error) => setError(err.message));

  const handleAdd = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    setIsAdding(true);
    run(createCategory(newName.trim()).then(() => setNewName(""))).finally(() =>
      setIsAdding(false)
    );
  };

  const handleRename = (id: number) => {
    run(renameCategory(id, editName.trim()).then(() => setEditingId(null)));
  };

  const handleDelete = (category: Category) => {
    if (window.confirm(`Xóa danh mục "${category.name}"?`)) {
      run(deleteCategory(category.id));
    }
  };

  const rows = categories.map((category) => (
    <tr key={category.id} className="border-t border-gray-100">
      <td className="py-3 pr-3 text-gray-500 w-16">#{category.id}</td>
      <td className="py-3 pr-3">
        {editingId === category.id ? (
          <form
            className="flex gap-2"
            onSubmit={(e) => {
              e.preventDefault();
              handleRename(category.id);
            }}
          >
            <input
              autoFocus
              required
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
              onKeyDown={(e) => e.key === "Escape" && setEditingId(null)}
              className="border border-gray-300 px-2 py-1 rounded-lg w-full max-w-xs"
              aria-label="Tên danh mục"
            />
            <button type="submit" className="text-[#06923E] font-medium cursor-pointer hover:underline">
              Lưu
            </button>
            <button
              type="button"
              onClick={() => setEditingId(null)}
              className="text-gray-500 cursor-pointer hover:underline"
            >
              Hủy
            </button>
          </form>
        ) : (
          <span className="font-medium">{category.name}</span>
        )}
      </td>
      <td className="py-3 text-right whitespace-nowrap">
        {editingId !== category.id && (
          <div className="flex gap-3 justify-end">
            <button
              type="button"
              onClick={() => {
                setEditingId(category.id);
                setEditName(category.name);
              }}
              className="text-[#FF7601] font-medium cursor-pointer hover:underline"
            >
              Đổi tên
            </button>
            {/* The gateway only lets the owner delete categories. */}
            {user.roleId === "owner" && (
              <button
                type="button"
                onClick={() => handleDelete(category)}
                className="text-red-600 font-medium cursor-pointer hover:underline"
              >
                Xóa
              </button>
            )}
          </div>
        )}
      </td>
    </tr>
  ));

  return (
    <div className="flex flex-col gap-4">
      <form onSubmit={handleAdd} className="flex gap-3">
        <input
          required
          placeholder="Tên danh mục mới"
          value={newName}
          onChange={(e) => setNewName(e.target.value)}
          className="border border-gray-300 p-2 rounded-lg bg-white text-sm w-full md:w-72"
          aria-label="Tên danh mục mới"
        />
        <button
          type="submit"
          disabled={isAdding}
          className="px-5 h-10 rounded-4xl bg-[#FF7601] text-white cursor-pointer hover:bg-[#ff8f33] disabled:bg-orange-300 whitespace-nowrap"
        >
          + Thêm
        </button>
      </form>
      <p className="text-xs text-gray-500">
        Đổi tên danh mục cũng cập nhật tên trên mọi sản phẩm thuộc danh mục. Chỉ xóa được danh mục
        không còn sản phẩm.
      </p>

      {error && <p className="bg-red-50 text-red-700 text-sm rounded-lg px-4 py-2">{error}</p>}

      <div className="bg-white rounded-2xl shadow p-5 overflow-x-auto">
        {isLoading ? (
          <div className="flex justify-center p-10">
            <AiOutlineLoading3Quarters className="animate-spin" size={32} />
          </div>
        ) : categories.length === 0 ? (
          <p className="text-center text-gray-500 p-10">Chưa có danh mục nào.</p>
        ) : (
          <table className="w-full text-sm">
            <tbody>{rows}</tbody>
          </table>
        )}
      </div>
    </div>
  );
};

export default CategoriesTab;
