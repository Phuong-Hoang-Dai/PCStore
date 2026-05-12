import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Home from "./pages/Home";
import Manage from "./pages/Manage";
import MainLayout from "./layouts/MainLayout";
import RequreAdmin from "./features/auth/RequreAdmin";
import Login from "./features/admin/Login";
function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<MainLayout />}>
          <Route index element={<Home />}></Route>
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
