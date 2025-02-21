// frontendnuevo/src/pages/HomePage.jsx
import React from "react";
import { useNavigate } from "react-router-dom";
import Header from "../components/Header";
import NavigationBar from "../components/NavigationBar";
import foto1 from "../images/foto1.png"; // Importa la imagen

const HomePage = () => {
  return (
    <div className="flex flex-col min-h-screen bg-white">
      <Header />
      <NavigationBar />
      {/* Agrega la imagen aquí */}
      <div className="w-full flex justify-center py-4 bg-white">
        <img
          src={foto1}
          alt="Foto de bienvenida"
          className="max-w-full h-auto"
        />
      </div>
      <main className="flex-grow container mx-auto px-4 py-4 text-center">
        <h1 className="text-3xl font-bold text-gray-900 mb-4 md:text-4xl">
          ¡Bienvenido a PrediApp!
        </h1>
        <p className="text-base text-gray-700 mb-6 md:text-lg">
          Predice los resultados de la Fórmula 1 y compite con tus amigos.
        </p>
        <button className="bg-red-500 text-white font-bold py-3 px-6 rounded-lg hover:bg-red-600 transition duration-300">
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
