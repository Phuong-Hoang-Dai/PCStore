import { apiFetch } from "../../utils/api";

// Through the gateway: the catalog service isn't reachable from the browser in Docker.
export const fetchProducts = async () => {
  try {
    const response = await apiFetch("/product?limit=10&offset=0");
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