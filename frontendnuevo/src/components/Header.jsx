// frontendnuevo/src/components/Header.jsx
import React, { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { FaUserCircle } from "react-icons/fa"; // Ya instalado con `npm install react-icons`

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const navigate = useNavigate();

  const toggleMenu = () => setIsMenuOpen(!isMenuOpen);

  return (
    <header className="bg-white text-black py-3 sticky top-0 z-50 shadow-md rounded-t-lg">
      <div className="container mx-auto px-4 flex items-center justify-between">
        {/* Botón hamburguesa (visible solo en móviles) */}
        <button
          className="md:hidden text-gray-800 hover:text-red-500 focus:outline-none"
          onClick={toggleMenu}
        >
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M4 6h16M4 12h16m-7 6h7"
            />
          </svg>
        </button>

        {/* Logo "predi" */}
        <Link
          to="/"
          className="text-xl font-bold text-red-500 hover:text-red-600 transition-colors duration-300 cursive"
        >
          predi
        </Link>

        {/* Ícono circular punteado (avatar, visible solo en escritorio) */}
        <div className="md:flex items-center space-x-2 hidden">
          <svg
            className="w-6 h-6 text-gray-600 hover:text-gray-800"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
            xmlns="http://www.w3.org/2000/svg"
          >
            <circle
              cx="12"
              cy="12"
              r="10"
              strokeWidth="2"
              strokeDasharray="5 5"
            />
          </svg>
        </div>

        {/* Menú de navegación (oculto en móviles por defecto, visible en md) */}
        <nav
          className={`${
            isMenuOpen ? "flex" : "hidden"
          } md:flex flex-col md:flex-row absolute md:static top-12 left-0 w-full md:w-auto bg-white md:bg-transparent p-4 md:p-0 space-y-4 md:space-y-0 md:space-x-6`}
        >
          <Link
            to="/"
            className="text-gray-800 hover:text-red-500 transition-colors duration-200 no-underline"
            onClick={() => setIsMenuOpen(false)}
          >
            inicio
          </Link>
          <Link
            to="/prodesPage"
            className="text-gray-800 hover:text-red-500 transition-colors duration-200 no-underline"
            onClick={() => setIsMenuOpen(false)}
          >
            prónosticos
          </Link>
          <Link
            to="/groups"
            className="text-gray-800 hover:text-red-500 transition-colors duration-200 no-underline"
            onClick={() => setIsMenuOpen(false)}
          >
            grupos
          </Link>
          <Link
            to="/results"
            className="text-gray-800 hover:text-red-500 transition-colors duration-200 no-underline"
            onClick={() => setIsMenuOpen(false)}
          >
            resultados
          </Link>
        </nav>
      </div>
    </header>
  );
};

export default Header;
