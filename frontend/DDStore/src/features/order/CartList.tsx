import { useDispatch, useSelector } from "react-redux";
import { changeByAmout, increment, selectCart } from "./cartSlice";
import { Swiper, SwiperSlide } from "swiper/react";
import { Scrollbar } from "swiper/modules";
import { IoChevronDownSharp, IoChevronUpSharp } from "react-icons/io5";
import { IoMdClose } from "react-icons/io";

const CartList = () => {
  const cartItem = useSelector(selectCart);
  const dispatch = useDispatch();
  const renderItems = cartItem.order.map((item) => (
    <SwiperSlide>
      <div className="flex flex-row items-center justify-around text-sm h-30 font-medium">
        <img
          src="https://cdn.hstatic.net/products/1000288298/11498_dsc05320_47c7abb602c949308c1f2bc50c3af657_master.jpg"
          className="w-2/8 h-auto"
        ></img>
        <div className="w-4/8 flex flex-col justify-between h-1/2 text-start px-2">
          <span>{item.name}</span>
          <span className="font-light ">{item.description}</span>
          <span className="text-end">{item.price}</span>
        </div>
        <span className="w-1/8 pr-3 flex flex-col items-center">
          <IoChevronUpSharp
            className="cursor-pointer"
            onClick={() => {
              dispatch(increment(item));
            }}
          />
          <span>{item.quantityOrder}</span>
          <IoChevronDownSharp
            className="cursor-pointer"
            onClick={() => {
              dispatch(increment(item));
            }}
          />
        </span>
        <IoMdClose
          className="w-1/8 cursor-pointer"
          onClick={() =>
            dispatch(changeByAmout({ item, quantity: item.quantityOrder }))
          }
        />
      </div>
    </SwiperSlide>
  ));

  const total: number = cartItem.order.reduce((sum: number, item): number => {
    return sum + item.price * item.quantityOrder;
  }, 0);
  return (
    <div className="bg-white px-8 rounded-xl relative shadow-lg w-screen md:w-md z-0 h-145 ">
      <span className=" w-0 h-0 border-[10px] -z-10 -top-[20px] right-2 border-b-white border-t-transparent border-x-transparent  absolute"></span>
      <span className=" w-[20px] h-0 border-b-[1px]  -z-10 border-solid top-0 right-20 border-white absolute"></span>
      <div className="flex w-full items-center justify-center pt-3">
        <div className="text-2xl font-medium border-neutral-400 border-b-1 w-5/6 my-2 pb-2">
          Giỏ hàng
        </div>
      </div>
      <Swiper
        modules={[Scrollbar]}
        scrollbar={{ draggable: true }}
        direction="vertical"
        className="max-h-90 h-fit"
        slidesPerView="auto"
      >
        {renderItems}
      </Swiper>
      <div className="w-full flex justify-between items-center text-2xl mt-3">
        <span>Tổng tiền:</span>
        <span>{total}</span>
      </div>
      <div className="w-full text-[20px] p-2 mt-9 bg-red-600 text-white rounded-4xl shadow-2xl cursor-pointer hover:bg-red-400">
        Hoàn tất đơn hàng
      </div>
    </div>
  );
};

export default CartList;
