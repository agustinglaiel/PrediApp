// frontendnuevo/src/pages/HomePage.jsx
import React from "react";
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";

const HomePage = () => {
  const navigate = useNavigate();

  return (
    <div className="flex flex-col min-h-screen bg-white">
      <Header />
      <main className="flex-grow container mx-auto px-4 py-6 text-center">
        <h1 className="text-2xl font-bold text-gray-900 mb-4 md:text-4xl">
          ¡Bienvenido a PrediApp!
        </h1>
        <p className="text-base text-gray-700 mb-6 md:text-lg">
          Predice los resultados de la Fórmula 1 y compite con tus amigos.
        </p>
        <button
          onClick={() => navigate("/prodesPage")}
          className="bg-red-500 text-white font-bold py-2 px-4 rounded-lg hover:bg-red-600 transition-colors duration-300 md:py-2 md:px-6"
        >
          ¡Haz tu Pronóstico!
        </button>
      </main>
      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm md:text-base">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default HomePage;
