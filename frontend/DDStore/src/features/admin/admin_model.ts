
export type Category = {
  id: number;
  name: string;
}

// A product as the catalog returns it (the storefront maps it to a smaller Product).
export type AdminProduct = {
  id: string;
  name: string;
  description: string;
  brand: string;
  cate: Category;
  price: number;
  type: string;
  state: number;
  images: { url: string; role: string }[] | null;
}

export type ProductInput = {
  name: string;
  description: string;
  brand: string;
  categoryId: number;
  price: number;
  imageUrl: string;
}

export type Stock = {
  productId: string;
  quantity: number;
}

export const toProductInput = (p: AdminProduct): ProductInput => ({
  name: p.name,
  description: p.description,
  brand: p.brand,
  categoryId: p.cate?.id ?? 0,
  price: p.price,
  imageUrl: p.images?.[0]?.url ?? "",
});
