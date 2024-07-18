import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import SignUp from "./components/Signup";
import Login from "./components/Login";
import Home from "./components/Home";
import AdminPanel from "./components/AdminPanel";
import UpdateUser from "./components/UpdateUser";
import "./styles/App.css";

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="/login" element={<Login />} />
        <Route path="/signup" element={<SignUp />} />
        <Route path="/admin" element={<AdminPanel />} />
        <Route path="/update-user/:id" element={<UpdateUser />} />
      </Routes>
    </Router>
  );
};

export default App;
