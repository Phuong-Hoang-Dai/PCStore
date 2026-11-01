
export const FALLBACK_IMAGE =
  "https://cdn.hstatic.net/products/1000288298/11498_dsc05320_47c7abb602c949308c1f2bc50c3af657_master.jpg";

export const ProductState = {
  IsAvailable: 0,
  OutOfStock: 1,
} as const;

export type Product = {
  id: string;
  name: string;
  description: string;
  price: number;
  image: string;
  state: number;
  quantityOrder: number;
}

type CatalogProduct = {
  id: string;
  name: string;
  description: string;
  price: number;
  state: number;
  images?: { url: string; role: string }[] | null;
}

export const mappingProduct = (data: Array<CatalogProduct> | null) => {
  const productList: Product[] = [];
  for (const key of data ?? []) {
    const product: Product = {
      id: key["id"],
      name: key["name"],
      description: key["description"],
      price: key["price"],
      image: key["images"]?.[0]?.url || FALLBACK_IMAGE,
      state: key["state"],
      quantityOrder: 0,
    };
    productList.push(product);
  }
  return productList;
}
