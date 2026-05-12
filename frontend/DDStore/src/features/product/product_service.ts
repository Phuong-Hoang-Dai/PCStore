const url = import.meta.env.VITE_API_URL + "/product"


export const fetchProducts = async () => {
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