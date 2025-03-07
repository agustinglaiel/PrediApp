import React from "react";
import LoginForm from "../components/LoginForm";
import Header from "../components/Header";

const LoginPage = () => {
  return (
    <div className="flex justify-center items-center min-h-screen bg-gray-50">
      <LoginForm />
      <footer className="bg-gray-200 text-gray-700 text-center py-3 text-sm w-full fixed bottom-0 left-0">
        <p>© 2025 PrediApp</p>
      </footer>
    </div>
  );
};

export default LoginPage;
