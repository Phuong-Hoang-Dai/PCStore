const url = import.meta.env.VITE_CATALOG_URL + "/product"


export const fetchProducts = async () => {
  console.log(url + "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
  try {
    const response = await fetch(url + "?limit=10&offset=0");
    if (!response.ok) {
      throw new Error("Network response was not ok");
    }
    const data = await response.json();
    return data;
  } catch (error) {
    console.error("Failed to fetch products:", error);
    throw error;
  }
}