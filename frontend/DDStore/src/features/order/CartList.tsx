import { useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useNavigate } from "react-router-dom";
import { changeByAmout, decrement, increment, selectCart } from "./cartSlice";
import { Swiper, SwiperSlide } from "swiper/react";
import { Mousewheel, Scrollbar } from "swiper/modules";
import { IoChevronDownSharp, IoChevronUpSharp } from "react-icons/io5";
import { IoMdClose } from "react-icons/io";
import { AiOutlineLoading3Quarters } from "react-icons/ai";
import { selectUser } from "../auth/userReducer";
import { createCheckoutSession, createOrder } from "./order_service";
import { formatVND } from "../../utils/format";
import { newTransactionId } from "../../utils/api";

const CartList = ({ onClose }: { onClose?: () => void }) => {
  const cartItem = useSelector(selectCart);
  const user = useSelector(selectUser);
  const dispatch = useDispatch();
  const navigate = useNavigate();
  const [isPending, setIsPending] = useState(false);
  const [error, setError] = useState("");

  const isLoggedIn = user.id !== 0;
  const isEmpty = cartItem.order.length === 0;

  const handleCheckout = () => {
    if (isPending || isEmpty) {
      return;
    }
    if (!isLoggedIn) {
      // Guests order with just a name and an email.
      onClose?.();
      navigate("/checkout");
      return;
    }

    setIsPending(true);
    setError("");
    // One transaction id for the whole checkout, so both requests are logged under it.
    const transactionId = newTransactionId();
    createOrder(
      cartItem.order.map((item) => ({
        productId: item.id,
        quantity: item.quantityOrder,
      })),
      transactionId
    )
      .then((orderId) => createCheckoutSession(orderId, transactionId))
      .then((url) => {
        // The cart is cleared on the success page, so it survives a cancelled payment.
        window.location.assign(url);
      })
      .catch((err: Error) => {
        setIsPending(false);
        setError(err.message || "Không thể tạo đơn hàng, vui lòng thử lại");
      });
  };

  const renderItems = cartItem.order.map((item) => (
    <SwiperSlide className="h-30" key={item.id}>
      <div className="flex flex-row items-center justify-around text-sm h-30 font-medium">
        <img
          src={item.image}
          alt={item.name}
          className="w-2/8 aspect-square object-cover"
        ></img>
        <div className="w-4/8 flex flex-col justify-between h-1/2 text-start px-2">
          <span className="truncate">{item.name}</span>
          <span className="font-light truncate">{item.description}</span>
          <span className="text-end">{formatVND(item.price)}</span>
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
              dispatch(decrement(item));
            }}
          />
        </span>
        <IoMdClose
          className="w-1/8 cursor-pointer"
          onClick={() =>
            dispatch(changeByAmout({ item, quantity: -item.quantityOrder }))
          }
        />
      </div>
    </SwiperSlide>
  ));

  const total: number = cartItem.order.reduce((sum: number, item): number => {
    return sum + item.price * item.quantityOrder;
  }, 0);
  return (
    <div className="bg-white px-8 pb-6 rounded-xl relative shadow-lg w-screen md:w-md z-0 min-h-145 text-base">
      <span className=" w-0 h-0 border-[10px] -z-10 -top-[20px] right-2 border-b-white border-t-transparent border-x-transparent  absolute"></span>
      <span className=" w-[20px] h-0 border-b-[1px]  -z-10 border-solid top-0 right-20 border-white absolute"></span>
      <div className="flex w-full items-center justify-center pt-3">
        <div className="text-2xl font-medium border-neutral-400 border-b-1 w-5/6 my-2 pb-2">
          Giỏ hàng
        </div>
      </div>
      {isEmpty ? (
        <div className="h-90 flex items-center justify-center text-gray-500">
          Giỏ hàng trống
        </div>
      ) : (
        <Swiper
          modules={[Scrollbar, Mousewheel]}
          scrollbar={{ draggable: true }}
          mousewheel={true}
          direction="vertical"
          className="h-90"
          slidesPerView="auto"
          spaceBetween={4}
        >
          {renderItems}
        </Swiper>
      )}
      <div className="w-full flex justify-between items-center text-2xl mt-3">
        <span>Tổng tiền:</span>
        <span>{formatVND(total)}</span>
      </div>
      {error && <p className="text-red-500 text-sm mt-3">{error}</p>}
      <button
        type="button"
        onClick={handleCheckout}
        disabled={isPending || isEmpty}
        className="w-full h-12 text-[20px] p-2 mt-6 bg-red-600 text-white rounded-4xl shadow-2xl cursor-pointer hover:bg-red-400 disabled:bg-red-300 disabled:cursor-not-allowed flex justify-center items-center"
      >
        {isPending ? (
          <AiOutlineLoading3Quarters className="animate-spin" />
        ) : (
          "Hoàn tất đơn hàng"
        )}
      </button>
    </div>
  );
};

export default CartList;
