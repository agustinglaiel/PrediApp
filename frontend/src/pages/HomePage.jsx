import React from "react";
import "../styles/HomePage.css"; // Importar el nuevo archivo de estilos
import { useNavigate } from "react-router-dom";
const HomePage = () => {
  const navigate = useNavigate();
  return (
    <div className="home-container">
      <header className="header">
        <div className="logo">
          {/* Reemplaza esto con tu logo real */}
          <img src="/logo.svg" alt="PrediApp Logo" />
        </div>
        <nav className="nav">
          <a href="/">Inicio</a>
          <a href="/prodesPage">Próximos Eventos</a>
          <a href="/results">Resultados</a>
          <a href="/login">Login</a>
          <a href="/signup">Signup</a>
        </nav>
      </header>
      <section className="hero">
        <div className="hero-content">
          <h1 className="hero-title">¡Bienvenido a PrediApp!</h1>
          <p className="hero-subtitle">
            Predice los resultados de la Fórmula 1 y compite con tus amigos.
          </p>
          <button className="hero-cta" onClick={() => navigate("/prodesPage")}>
            ¡Haz tu Pronóstico!
          </button>
        </div>
        <div className="hero-image">
          <img src="f1-car.jpg" alt="F1 Car" />
        </div>
      </section>
      <section className="featured-events">
        <h2>Próximos Eventos Destacados</h2>
        {/* Aquí irán los componentes de los próximos eventos */}
      </section>
      <section className="latest-results">
        <h2>Últimos Resultados</h2>
        {/* Aquí irán el componente del carrusel de resultados */}
      </section>
      <footer className="footer">
        <p>© 2024 PrediApp</p>
      </footer>
    </div>
  );
};
export default HomePage;
