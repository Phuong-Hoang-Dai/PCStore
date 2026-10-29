import type { Product } from "./product_model";
import UpdateQuantityOrder from "./UpdateQuantityOrder";
import { formatVND } from "../../utils/format";

const ProductItem = ({ product }: { product: Product }) => {
  return (
    <div
      key={product.id}
      className="pb-5 gap-5 flex flex-col overflow-hidden px-1 w-full rounded-2xl shadow-lg bg-white"
    >
      <img
        src={product.image}
        alt={product.name}
        className="rounded-t-2xl w-full aspect-square object-cover"
      />
      <div className="px-3 flex flex-col gap-2 text-sm">
        <h2 className="font-medium uppercase truncate">{product.name}</h2>
        <span className="text-red-600 font-medium">
          {formatVND(product.price)}
        </span>
        <UpdateQuantityOrder product={product} />
      </div>
    </div>
  );
};

export default ProductItem;
