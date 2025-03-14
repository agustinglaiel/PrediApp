import React, { useState, useEffect } from "react";

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [showHamburger, setShowHamburger] = useState(false);

  const toggleMenu = () => {
    setIsMenuOpen((prev) => !prev);
  };

  useEffect(() => {
    const handleScroll = () => {
      const currentScrollPos = window.scrollY;
      setShowHamburger(currentScrollPos > 64); // Ajusta si la altura del Header cambia
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  const Link = ({ to, children, className, onClick }) => (
    <button className={className} onClick={onClick}>
      {children}
    </button>
  );

  return (
    <div className="relative">
      <header className="bg-red-600 text-white w-full z-50 shadow-md fixed top-0 left-0">
        <div className="container mx-auto px-4 py-4 flex items-center justify-between">
          {/* Espacio reservado para la hamburguesa (left) - visible solo en mobile */}
          <div className="w-6 md:hidden">
            {showHamburger && (
              <button
                className="text-white hover:text-gray-200 focus:outline-none"
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
            )}
          </div>

          {/* Logo/Texto "PREDI" (center) - ahora con flex para mantener el centrado */}
          <div className="flex-1 flex justify-center">
            <Link
              to="/"
              className="text-2xl md:text-3xl font-bold tracking-wide"
            >
              PREDI
            </Link>
          </div>

          {/* Espacio reservado para el logo (right) - visible solo en escritorio */}
          <div className="w-6 md:w-auto">
            <div className="hidden md:block">
              <img
                src="/api/placeholder/32/32"
                alt="Logo"
                className="w-8 h-8 rounded-full border-2 border-white"
              />
            </div>
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
