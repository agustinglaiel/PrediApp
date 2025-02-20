import React from "react";
import { Link, useNavigate } from "react-router-dom";
import "../styles/components.css";
const Header = () => {
  const navigate = useNavigate();
  return (
    <header className="bg-gray-900 text-white py-6">
      <div className="container mx-auto flex items-center justify-between">
        {/* Logo */}
        <div className="logo">
          <Link
            to="/"
            className="text-2xl font-bold text-red-500 hover:text-red-400 transition-colors duration-300"
          >
            PrediApp
          </Link>
        </div>
        {/* Navigation */}
        <nav className="space-x-6">
          <Link
            to="/"
            className="hover:text-gray-300 transition-colors duration-200"
          >
            Inicio
          </Link>
          <Link
            to="/prodesPage"
            className="hover:text-gray-300 transition-colors duration-200"
          >
            Próximos Eventos
          </Link>
          <Link
            to="/results"
            className="hover:text-gray-300 transition-colors duration-200"
          >
            Resultados
          </Link>
          <Link
            to="/login"
            className="hover:text-gray-300 transition-colors duration-200"
          >
            Login
          </Link>
          <Link
            to="/signup"
            className="hover:text-gray-300 transition-colors duration-200"
          >
            Signup
          </Link>
        </nav>
      </div>
    </header>
  );
};
export default Header;
