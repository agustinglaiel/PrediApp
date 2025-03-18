import React from "react";

const SignOutAlert = ({ isOpen, onClose, onConfirm }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white p-6 rounded-lg shadow-lg max-w-md w-full mx-4 md:mx-8 lg:mx-16">
        <p className="text-center text-gray-800 mb-4">
          ¿Estás seguro de que quieres cerrar sesión?
        </p>
        <div className="flex justify-center space-x-4">
          <button
            onClick={onClose}
            className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600 focus:outline-none"
          >
            Cancelar
          </button>
          <button
            onClick={onConfirm}
            className="bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600 focus:outline-none"
          >
            Confirmar
          </button>
        </div>
      </div>
    </div>
  );
};

export default SignOutAlert;
