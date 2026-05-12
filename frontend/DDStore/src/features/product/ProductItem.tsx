import type { Product } from "./product_model";
import UpdateQuantityOrder from "./UpdateQuantityOrder";

const ProductItem = ({ product }: { product: Product }) => {
  return (
    <div
      key={product.id}
      className="pb-5 gap-5 flex flex-col overflow-hidden px-1 w-full rounded-2xl shadow-lg bg-white"
    >
      <img
        src="https://cdn.hstatic.net/products/1000288298/11498_dsc05320_47c7abb602c949308c1f2bc50c3af657_master.jpg"
        alt={product.name}
        className="rounded-t-2xl-2xl i w-full h-auto"
      />
      <div className="px-3 flex flex-col gap-2 text-sm">
        <h2 className="font-medium uppercase truncate">{product.name}</h2>
        <span className="text-red-600 font-medium">{product.price}đ</span>
        <UpdateQuantityOrder product={product} />
      </div>
    </div>
  );
};

export default ProductItem;
