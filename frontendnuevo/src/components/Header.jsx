import React, { useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import SignOutAlert from "./SignOutAlert";

const Header = () => {
  const [isMenuOpen, setIsMenuOpen] = useState(false);
  const [showHamburger, setShowHamburger] = useState(false);
  const [showSignOutModal, setShowSignOutModal] = useState(false);
  const navigate = useNavigate();

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

  // Verificar si el usuario está autenticado
  const isAuthenticated = !!localStorage.getItem("jwtToken");

  // Función para manejar el clic en el logo
  const handleLogoClick = () => {
    if (isAuthenticated) {
      setShowSignOutModal(true); // Mostrar el modal si está autenticado
    } else {
      navigate("/login"); // Redirigir a login si no está autenticado
    }
  };

  // Función para confirmar el cierre de sesión
  const confirmSignOut = () => {
    localStorage.removeItem("jwtToken");
    localStorage.removeItem("userId");
    setShowSignOutModal(false);
    navigate("/");
  };

  // Función para cerrar el modal sin acción
  const closeSignOutModal = () => {
    setShowSignOutModal(false);
  };

  const Link = ({ to, children, className, onClick }) => (
    <button className={className} onClick={onClick}>
      {children}
    </button>
  );

  return (
    <div className="relative">
      <header className="bg-red-600 text-white w-full z-50 shadow-md fixed top-0 left-0 mb-3">
        <div className="container mx-auto px-4 py-4 flex items-center">
          {/* Contenedor izquierdo para hamburguesa */}
          <div className="w-6 flex-none">
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

          {/* Logo/Texto "PREDI" (center) - usando absolute para centrarlo siempre */}
          <div className="absolute left-0 right-0 flex justify-center pointer-events-none">
            <Link
              to="/"
              className="text-2xl md:text-3xl font-bold tracking-wide pointer-events-auto"
              onClick={() => navigate("/")}
            >
              PREDI
            </Link>
          </div>

          {/* Espacio para mantener el flex balanceado */}
          <div className="flex-1"></div>

          {/* Logo como botón de Log In/Sign Out (right) */}
          <div className="flex-none mr-2 md:mr-6 flex items-center">
            {" "}
            {/* Añadido flex items-center */}
            <button
              onClick={handleLogoClick}
              className="focus:outline-none flex items-center justify-center" /* Añadido flex items-center justify-center */
              title={isAuthenticated ? "Sign Out" : "Log In"}
            >
              <div className="w-9 h-9 rounded-full overflow-hidden">
                <img
                  src="/images/user-logo.jpg"
                  alt={isAuthenticated ? "Sign Out" : "Log In"}
                  className="w-full h-full object-contain"
                />
              </div>
            </button>
          </div>
        </div>

        {/* Menú móvil (hamburger) */}
        {isMenuOpen && (
          <nav className="md:hidden bg-red-600 px-4 pb-4 flex flex-col space-y-2 text-center shadow-md">
            <Link
              to="/"
              className="hover:text-gray-200"
              onClick={() => {
                setIsMenuOpen(false);
                navigate("/");
              }}
            >
              Inicio
            </Link>
            <Link
              to="/pronosticos"
              className="hover:text-gray-200"
              onClick={() => {
                setIsMenuOpen(false);
                navigate("/pronosticos");
              }}
            >
              Prónosticos
            </Link>
            <Link
              to="/grupos"
              className="hover:text-gray-200"
              onClick={() => {
                setIsMenuOpen(false);
                navigate("/grupos");
              }}
            >
              Grupos
            </Link>
            <Link
              to="/resultados"
              className="hover:text-gray-200"
              onClick={() => {
                setIsMenuOpen(false);
                navigate("/resultados");
              }}
            >
              Resultados
            </Link>
          </nav>
        )}
      </header>

      {/* Modal de confirmación de Sign Out */}
      <SignOutAlert
        isOpen={showSignOutModal}
        onClose={closeSignOutModal}
        onConfirm={confirmSignOut}
      />
    </div>
  );
};

export default Header;
