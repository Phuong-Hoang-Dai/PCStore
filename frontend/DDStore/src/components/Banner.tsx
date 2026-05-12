import { useRef } from "react";
import { FaArrowLeft, FaArrowRight } from "react-icons/fa6";
import { A11y, Autoplay, Navigation, Pagination } from "swiper/modules";
import { Swiper, SwiperSlide } from "swiper/react";

const Banner = () => {
  const prevRef = useRef<HTMLButtonElement | null>(null);
  const nextRef = useRef<HTMLButtonElement | null>(null);
  console.log("re-render : ref: " + prevRef);
  const images = [
    "https://theme.hstatic.net/1000288298/1001020793/14/slide_1_img.jpg?v=1540",
    "https://theme.hstatic.net/1000288298/1001020793/14/slide_2_img.jpg?v=1540",
    "https://theme.hstatic.net/1000288298/1001020793/14/slide_4_img.jpg?v=1540",
  ];

  const renderBanner = images.map((src, i) => (
    <SwiperSlide key={i}>
      <img src={src} alt={"Banner " + i} />
    </SwiperSlide>
  ));

  return (
    <>
      <div className="relative">
        <button
          ref={prevRef}
          className="absolute top-1/2 left-4 -translate-y-1/2 z-10 bg-white p-2 rounded-full shadow hover:bg-gray-200 transition"
        >
          <FaArrowLeft className="text-xl text-black" />
        </button>
        <button
          ref={nextRef}
          className="absolute top-1/2 right-4 -translate-y-1/2 z-10 bg-white p-2 rounded-full shadow hover:bg-gray-200 transition"
        >
          <FaArrowRight className="text-xl text-black" />
        </button>
        <Swiper
          className=""
          modules={[Navigation, Pagination, Autoplay, A11y]}
          slidesPerView={1}
          navigation={{
            nextEl: nextRef.current,
            prevEl: prevRef.current,
          }}
          onInit={(swiper) => {
            if (
              typeof swiper.params.navigation !== "boolean" &&
              swiper.params.navigation
            ) {
              swiper.params.navigation.prevEl = prevRef.current;
              swiper.params.navigation.nextEl = nextRef.current;
            }
            swiper.navigation.init();
            swiper.navigation.update();
          }}
          pagination={{
            clickable: false,
          }}
          speed={1500}
          autoplay={{ delay: 3000 }}
          loop
        >
          {renderBanner}
        </Swiper>
      </div>
    </>
  );
};

export default Banner;
