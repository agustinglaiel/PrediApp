// frontendnuevo/src/App.jsx
import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import ProdeSessionPage from "./pages/ProdeSessionPage";

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route
          path="/pronosticos/:session_id"
          element={<ProdeSessionPage />}
        />{" "}
        {/* Nueva ruta para pronósticos */}
      </Routes>
    </Router>
  );
};

export default App;
