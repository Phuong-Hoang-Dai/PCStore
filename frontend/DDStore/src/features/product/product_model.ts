
export type Product = {
  id: number;
  name: string;
  description: string;
  price: number;
  quantity: number;
  quantityOrder: number;
}

export const mappingProduct = (data: Array<Product>)  => {
  let productList: Product[] = [];
  for (const key of data) {
    const product: Product = {
      id: key["id"],
      name: key["name"],
      description: key["description"],
      price: key["price"],
      quantity:key["quantity"],
      quantityOrder:0,
    };
    productList.push(product);
  }
  return productList;
}