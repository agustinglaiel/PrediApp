// pages/ProdeDispatcher.jsx
import React from "react";
import { useParams, useLocation } from "react-router-dom";
import ProdeSessionPage from "./ProdeSessionPage";
import ProdeRacePage from "./ProdeRacePage";

const ProdeDispatcher = () => {
  const { session_id } = useParams();
  const { state } = useLocation();

  // Si no vino state, podrías hacer un fetch en otro momento...
  const isRace = state?.sessionType === "Race" && state?.sessionName === "Race";

  if (isRace) {
    return <ProdeRacePage />;
  } else {
    return <ProdeSessionPage />;
  }
};

export default ProdeDispatcher;
