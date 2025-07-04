import { BrowserRouter, Routes, Route } from "react-router";
import { useEffect } from "react";
import Index from "./Index/index.jsx";
import Login from "./Login/login.jsx";
import SignUp from "./SignUp/signup.jsx";
import Home from "./Home/home.jsx";

import "./App.css";

const App = () => {
  useEffect(() => {
    document.body.style.visibility = "visible";
  }, []);

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Index />} />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/home" element={<Home />} />
        <Route path="/prehome" element={<Home />} />
      </Routes>
    </BrowserRouter>
  );
};

export default App;
