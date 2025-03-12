// frontendnuevo/src/App.jsx
import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";

import HomePage from "./pages/HomePage";
import LoginPage from "./pages/LoginPage";
import ProdeDispatcher from "./pages/ProdeDispatcher";
import PronosticosPage from "./pages/PronosticosPage";
import ProdeSessionResultPage from "./pages/ProdeSessionResultPage";
import ProdeRaceResultPage from "./pages/ProdeRaceResultPage";

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/login" element={<LoginPage />} />
        <Route path="/pronosticos" element={<PronosticosPage />} />
        <Route path="/pronosticos/:session_id" element={<ProdeDispatcher />} />
        <Route
          path="/pronosticos/result/:session_id"
          element={<ProdeSessionResultPage />}
        />
        <Route
          path="/pronosticos/result/race/:session_id"
          element={<ProdeRaceResultPage />}
        />
        {/* <Route path="*" element={<NotFoundPage />} /> */}
      </Routes>
    </Router>
  );
}

export default App;
