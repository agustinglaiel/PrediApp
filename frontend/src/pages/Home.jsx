import React from "react";
import "../styles/Home.css";
import { useNavigate } from "react-router-dom";

const Home = () => {
  const navigate = useNavigate();

  const handleClick = () => {
    navigate("/signup");
  };

  return (
    <div
      className="home-container"
      onClick={handleClick}
      style={{ cursor: "pointer" }}
    >
      <div className="welcome-box">Bienvenido</div>
    </div>
  );
};

export default Home;
