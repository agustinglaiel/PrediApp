import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import SignUp from "./pages/Signup";
import Login from "./pages/Login";
import ProdesPage from "./pages/ProdesPage";
import AdminPanel from "./pages/AdminPanel";
import UpdateUser from "./pages/UpdateUser";
import RaceWeekendPage from "./pages/RaceWeekendPage";
import { SessionsProvider } from "./contexts/SessionsContext"; // Importamos el proveedor del contexto
import RacePredictionForm from "./pages/RacePredictionForm";
import SessionPredictionForm from "./pages/SessionPredictionForm";
import Home from "./pages/Home";

const App = () => {
  return (
    <SessionsProvider>
      {" "}
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/prodesPage" element={<ProdesPage />} />
          <Route path="/login" element={<Login />} />
          <Route path="/signup" element={<SignUp />} />
          <Route path="/admin" element={<AdminPanel />} />
          <Route path="/update-user/:id" element={<UpdateUser />} />
          <Route path="/race-weekend" element={<RaceWeekendPage />} />
          <Route
            path="/prode/race/:sessionId"
            element={<RacePredictionForm />}
          />
          <Route
            path="/prode/session/:sessionId"
            element={<SessionPredictionForm />}
          />
        </Routes>
      </Router>
    </SessionsProvider>
  );
};

export default App;
