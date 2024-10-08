import React from "react";
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import SignUp from "./components/Signup";
import Login from "./components/Login";
import Home from "./components/Home";
import AdminPanel from "./components/AdminPanel";
import UpdateUser from "./components/UpdateUser";
import RaceWeekendPage from "./components/RaceWeekendPage";
import { SessionsProvider } from "./contexts/SessionsContext"; // Importamos el proveedor del contexto
import RacePredictionForm from "./components/RacePredictionForm";
import SessionPredictionForm from "./components/SessionPredictionForm";

const App = () => {
  return (
    <SessionsProvider>
      {" "}
      <Router>
        <Routes>
          <Route path="/" element={<Home />} />
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
