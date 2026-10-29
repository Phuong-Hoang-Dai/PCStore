import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import Manage from "./pages/Manage";
import OrderHistory from "./pages/OrderHistory";
import PaymentSuccess from "./pages/PaymentSuccess";
import PaymentCancel from "./pages/PaymentCancel";
import MainLayout from "./layouts/MainLayout";
import RequreAdmin from "./features/auth/RequreAdmin";
import Login from "./features/admin/Login";
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MainLayout />}>
          <Route index element={<Home />}></Route>
          <Route path="orders" element={<OrderHistory />} />
          <Route path="payment/success" element={<PaymentSuccess />} />
          <Route path="payment/cancel" element={<PaymentCancel />} />
        </Route>
        <Route path="/admin/login" element={<Login />} />
        <Route path="/manage" element={<RequreAdmin />}>
          <Route index element={<Manage />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
