import React from "react";
import "../styles/HomePage.css"; // Agrega esto
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";

const HomePage = () => {
  const navigate = useNavigate();
  return (
    <div>
      <Header />
      <main>
        <div>
          <div>
            <h1>¡Bienvenido a PrediApp!</h1>
            <p>
              Predice los resultados de la Fórmula 1 y compite con tus amigos.
            </p>
            <button onClick={() => navigate("/prodesPage")}>
              ¡Haz tu Pronóstico!
            </button>
          </div>

          <div>
            <h2>Próximos Eventos</h2>
            {/* Aquí va la lista de próximos eventos */}
          </div>

          <div>
            <h2>Últimos Resultados</h2>
            {/* Aquí va el carrusel de los últimos resultados */}
          </div>
        </div>
      </main>
      <footer>
        <p>© 2024 PrediApp</p>
      </footer>
    </div>
  );
};

export default HomePage;
