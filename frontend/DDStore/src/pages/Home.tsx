import Banner from "../components/Banner";
import ProductList from "../features/product/ProductList";

const Home = () => {
  return (
    <>
      <Banner />
      <div className=" mt-5">
        <ProductList />
      </div>
    </>
  );
};

export default Home;
