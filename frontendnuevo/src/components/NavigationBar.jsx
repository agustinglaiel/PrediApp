// src/components/NavigationBar.jsx
import React from "react";
import { useNavigate, useLocation } from "react-router-dom";

const NavigationBar = () => {
  const navigate = useNavigate();
  const location = useLocation();

  const navItems = [
    { path: "/", label: "inicio" },
    { path: "/pronosticos", label: "prónósticos" },
    { path: "/grupos", label: "grupos" },
    { path: "/resultados", label: "resultados" },
  ];

  return (
    <nav className="bg-gray-100 border-t border-gray-300">
      <div className="container mx-auto px-4 py-2 flex justify-around items-center">
        {navItems.map((item) => (
          <button
            key={item.path}
            onClick={() => navigate(item.path)}
            className={`px-4 py-2 text-sm font-medium rounded-md ${
              location.pathname === item.path
                ? "text-red-600 bg-gray-200"
                : "text-gray-700 hover:text-red-600"
            }`}
          >
            {item.label}
          </button>
        ))}
      </div>
    </nav>
  );
};

export default NavigationBar;
