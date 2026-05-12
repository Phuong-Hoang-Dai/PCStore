import { fetchProducts } from "./product_service";
import { useEffect, useState } from "react";
import { type Product, mappingProduct } from "./product_model";
import ProductItem from "./ProductItem";

const ProductList = () => {
  const [products, setProducts] = useState<Product[]>([]);
  useEffect(() => {
    fetchProducts().then((data) => {
      setProducts(mappingProduct(data["data"]));
    });
  }, []);

  const renderProducts = products.map((product) => (
    <ProductItem product={product} key={product.id} />
  ));

  return (
    <>
      <div className="grid grid-cols-2 md:grid-cols-3 gap-6 lg:grid-cols-5">
        {renderProducts}
      </div>
    </>
  );
};

export default ProductList;
