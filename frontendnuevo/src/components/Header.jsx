// src/components/Header.jsx
import React, { useState } from "react";

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);

  const toggleMenu = () => {
    setIsMenuOpen((prev) => !prev);
  };

  // Simulamos los Links como botones para la preview
  const Link = ({ to, children, className, onClick }) => (
    <button className={className} onClick={onClick}>
      {children}
    </button>
  );

  return (
    <div className="relative bg-gray-100">
      {" "}
      {/* Remove h-96 */}
      <header className="bg-red-600 text-white w-full z-50 shadow-md">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          {/* Botón hamburguesa (left) - visible solo en mobile */}
          <button
            className="md:hidden text-white hover:text-gray-200 focus:outline-none"
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
                d="M4 6h16M4 12h16M4 18h16"
              />
            </svg>
          </button>

          {/* Logo/Texto "PREDI" (center) */}
          <Link
            to="/"
            className="text-2xl md:text-3xl font-bold tracking-wide flex-1 text-center"
          >
            PREDI
          </Link>

          {/* Logo placeholder (right) - visible solo en escritorio */}
          <div className="hidden md:block">
            <img
              src="/api/placeholder/32/32"
              alt="Logo"
              className="w-8 h-8 rounded-full border-2 border-white"
            />
          </div>
        </div>

        {/* Menú móvil (hamburger) */}
        {isMenuOpen && (
          <nav className="md:hidden bg-red-600 px-4 pb-4 flex flex-col space-y-2 text-center shadow-md">
            <Link
              to="/"
              className="hover:text-gray-200"
              onClick={() => setIsMenuOpen(false)}
            >
              Inicio
            </Link>
            <Link
              to="/pronosticos"
              className="hover:text-gray-200"
              onClick={() => setIsMenuOpen(false)}
            >
              Prónosticos
            </Link>
            <Link
              to="/grupos"
              className="hover:text-gray-200"
              onClick={() => setIsMenuOpen(false)}
            >
              Grupos
            </Link>
            <Link
              to="/resultados"
              className="hover:text-gray-200"
              onClick={() => setIsMenuOpen(false)}
            >
              Resultados
            </Link>
          </nav>
        )}
      </header>
    </div>
  );
};

export default Header;
