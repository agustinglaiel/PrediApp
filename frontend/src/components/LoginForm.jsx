import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import { login } from "../api/users";

const LoginForm = ({ onLoginSuccess }) => {
  const [credentials, setCredentials] = useState({ email: "", password: "" });
  const [error, setError] = useState(null);
  const navigate = useNavigate();

  const handleChange = (e) => {
    setCredentials({ ...credentials, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    try {
      setError(null);
      const response = await login(credentials);
      if (onLoginSuccess) {
        onLoginSuccess(); // Llamamos al callback si se proporciona
      }
      navigate("/", { replace: true });
    } catch (err) {
      setError(
        err.message ||
          "Error al iniciar sesión. Verifica tus credenciales o contacta al soporte."
      );
    }
  };

  return (
    <div className="w-full max-w-sm p-8 bg-white rounded-2xl shadow-2xl border border-gray-100">
      <h2 className="text-3xl font-semibold text-center text-gray-800 mb-6">
        Iniciar Sesión
      </h2>
      {error && (
        <p className="text-red-600 text-center mb-4 text-sm font-medium">
          {error}
        </p>
      )}
      <form onSubmit={handleSubmit} className="space-y-5">
        <div>
          <label
            htmlFor="email"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Email
          </label>
          <input
            type="email"
            id="email"
            name="email"
            value={credentials.email}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all duration-200"
            placeholder="tu@email.com"
            required
            autoComplete="email" // Añadimos autocomplete para el campo de email
          />
        </div>
        <div>
          <label
            htmlFor="password"
            className="block text-sm font-medium text-gray-700 mb-1"
          >
            Contraseña
          </label>
          <input
            type="password"
            id="password"
            name="password"
            value={credentials.password}
            onChange={handleChange}
            className="w-full px-4 py-3 border border-gray-200 rounded-xl focus:outline-none focus:ring-2 focus:ring-red-500 focus:border-transparent transition-all duration-200"
            placeholder="••••••••"
            required
            autoComplete="current-password" // Añadimos autocomplete para el campo de contraseña
          />
        </div>
        <button
          type="submit"
          className="w-full bg-red-600 text-white py-3 px-4 rounded-xl hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 transition-all duration-200 font-medium"
        >
          Iniciar Sesión
        </button>
      </form>
    </div>
  );
};

export default LoginForm;
