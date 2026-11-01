// Offset paging: the APIs return no total, so "next" is enabled while a page comes back full.
const Pager = ({
  page,
  hasNext,
  onChange,
}: {
  page: number;
  hasNext: boolean;
  onChange: (page: number) => void;
}) => (
  <div className="flex items-center justify-end gap-3 text-sm mt-4">
    <button
      type="button"
      disabled={page === 0}
      onClick={() => onChange(page - 1)}
      className="px-4 h-9 rounded-4xl border border-gray-300 cursor-pointer hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed"
    >
      Trước
    </button>
    <span className="text-gray-600">Trang {page + 1}</span>
    <button
      type="button"
      disabled={!hasNext}
      onClick={() => onChange(page + 1)}
      className="px-4 h-9 rounded-4xl border border-gray-300 cursor-pointer hover:bg-gray-100 disabled:opacity-40 disabled:cursor-not-allowed"
    >
      Sau
    </button>
  </div>
);

export default Pager;
